package main

import (
	"context"

	"github.com/clls-dev/clls/pkg/clls"
	"github.com/pkg/errors"
	lsp "go.lsp.dev/protocol"
	"go.uber.org/zap"
)

type SemanticTokensOptions struct {
	lsp.WorkDoneProgressOptions
	Legend lsp.SemanticTokensLegend `json:"legend"`
	Range  *lsp.Range               `json:"range,omitempty"`
	Full   interface{}              `json:"full,omitempty"`
}

func (s *server) Initialize(context.Context, *lsp.InitializeParams) (*lsp.InitializeResult, error) {
	caps := lsp.ServerCapabilities{
		TextDocumentSync: lsp.TextDocumentSyncKindFull,
		SemanticTokensProvider: SemanticTokensOptions{
			Legend: clls.StandardSemanticTokensLegend,
			Full:   true,
		},
		DocumentFormattingProvider: true,
		RenameProvider:             true,
		DocumentHighlightProvider:  true,
	}
	s.l.Debug("server initialized", zap.Any("capabilities", caps))
	return &lsp.InitializeResult{
		ServerInfo: &lsp.ServerInfo{
			Name:    "clls",
			Version: "0.1.0",
		},
		Capabilities: caps}, nil

}

func (s *server) DidOpen(_ context.Context, params *lsp.DidOpenTextDocumentParams) error {
	docData := newDocumentData(params.TextDocument.Text)
	if pulled, ok := s.cache.pull(docData.contentHash); ok {
		s.openedDocs[params.TextDocument.URI] = pulled
		return nil
	}
	s.openedDocs[params.TextDocument.URI] = newDocumentData(params.TextDocument.Text)
	return nil
}

// This only supports full file changes
func (s *server) DidChange(_ context.Context, params *lsp.DidChangeTextDocumentParams) error {
	if len(params.ContentChanges) == 0 {
		return nil
	}

	docData := newDocumentData(params.ContentChanges[0].Text)

	if edd, ok := s.openedDocs[params.TextDocument.URI]; ok && edd.contentHash == docData.contentHash {
		return nil // data didn't change
	}

	odd, ok := s.openedDocs[params.TextDocument.URI]
	if !ok {
		return errors.New("document not opened")
	}

	delete(s.openedDocs, params.TextDocument.URI)
	s.cache.put(odd)

	if pulled, ok := s.cache.pull(docData.contentHash); ok {
		s.openedDocs[params.TextDocument.URI] = pulled
		return nil
	}

	s.openedDocs[params.TextDocument.URI] = docData
	return nil
}

func (s *server) DidClose(_ context.Context, params *lsp.DidCloseTextDocumentParams) error {
	dd, ok := s.openedDocs[params.TextDocument.URI]
	if !ok {
		return nil
	}
	delete(s.openedDocs, params.TextDocument.URI)
	s.cache.put(dd)
	return nil
}

func (s *server) Rename(_ context.Context, params *lsp.RenameParams) (*lsp.WorkspaceEdit, error) {
	p := params.Position
	line := int(p.Line)
	char := int(p.Character)
	newName := params.NewName

	var syms []*clls.Symbol
	if d, ok := s.openedDocs[params.TextDocument.URI]; ok && d.generatedSymbols {
		syms = d.symbols
	} else {
		mod, err := s.loadCLVM(params.TextDocument.URI)
		if err != nil {
			return nil, errors.Wrap(err, "parse module")
		}

		syms = mod.Symbols(s.l)
	}

	if d, ok := s.openedDocs[params.TextDocument.URI]; ok {
		d.symbols = syms
		d.generatedSymbols = true
	}

	for _, s := range syms {
		for _, st := range s.Tokens() {
			if st.Line != line {
				continue
			}
			if char < st.StartChar || st.EndChar() <= char {
				continue
			}

			edit := lsp.WorkspaceEdit{
				Changes: map[lsp.DocumentURI][]lsp.TextEdit{},
			}
			for _, t := range s.Tokens() {
				edit.Changes[t.DocumentURI] = append(edit.Changes[t.DocumentURI], lsp.TextEdit{
					Range:   t.Range(),
					NewText: newName,
				})
			}
			return &edit, nil
		}
	}

	return nil, errors.New("symbol not found")
}

func (s *server) Formatting(_ context.Context, params *lsp.DocumentFormattingParams) ([]lsp.TextEdit, error) {
	uriStr := params.TextDocument.URI

	fileStr, err := s.readFile(uriStr)
	if err != nil {
		return nil, errors.Wrap(err, "read file")
	}

	newText, linesCount, err := clls.Prettify(s.l, &params.Options, fileStr, uriStr)
	if err != nil {
		return nil, errors.Wrap(err, "prettify file content")
	}

	return []lsp.TextEdit{{
		Range: lsp.Range{
			Start: lsp.Position{Line: 0, Character: 0},
			End:   lsp.Position{Line: uint32(linesCount + 1), Character: 0},
		},
		NewText: newText,
	}}, nil
}

func (s *server) SemanticTokensFull(_ context.Context, params *lsp.SemanticTokensParams) (*lsp.SemanticTokens, error) {
	if d, ok := s.openedDocs[params.TextDocument.URI]; ok && d.generatedTokens {
		return &lsp.SemanticTokens{Data: d.semanticTokens}, nil
	}

	mod, err := s.loadCLVM(params.TextDocument.URI)
	if err != nil {
		return nil, errors.Wrap(err, "parse module")
	}

	data, err := mod.SemanticTokens(s.l)
	if err != nil {
		return nil, errors.Wrap(err, "semantic tokens from module")
	}

	s.l.Debug("generated semantic tokens", zap.Int("count", len(data)/5))

	if d, ok := s.openedDocs[params.TextDocument.URI]; ok {
		d.semanticTokens = data
		d.generatedTokens = true
	}

	return &lsp.SemanticTokens{Data: data}, nil
}

func (s *server) Shutdown(context.Context) error {
	s.down = true
	return nil
}

func (s *server) Exit(context.Context) error {
	s.exit = true
	return nil
}

func (s *server) Initialized(context.Context, *lsp.InitializedParams) error {
	return nil
}

func (s *server) DidSave(context.Context, *lsp.DidSaveTextDocumentParams) error {
	return nil // this is to gracefully ignore the event
}

func (s *server) DocumentHighlight(_ context.Context, params *lsp.DocumentHighlightParams) ([]lsp.DocumentHighlight, error) {
	p := params.Position
	line := int(p.Line)
	char := int(p.Character)

	var syms []*clls.Symbol
	if d, ok := s.openedDocs[params.TextDocument.URI]; ok && d.generatedSymbols {
		syms = d.symbols
	} else {
		mod, err := s.loadCLVM(params.TextDocument.URI)
		if err != nil {
			return nil, errors.Wrap(err, "parse module")
		}

		syms = mod.Symbols(s.l)
	}

	if d, ok := s.openedDocs[params.TextDocument.URI]; ok {
		d.symbols = syms
		d.generatedSymbols = true
	}

	for _, sym := range syms {
		for _, st := range sym.Tokens() {
			if st.Line != line {
				continue
			}
			if char < st.StartChar || st.EndChar() < char {
				continue
			}

			r := []lsp.DocumentHighlight(nil)
			for _, t := range sym.Tokens() {
				if t.DocumentURI != params.TextDocument.URI {
					continue
				}
				r = append(r, lsp.DocumentHighlight{Range: t.Range(), Kind: lsp.DocumentHighlightKindText})
			}
			s.l.Debug("will send", zap.Any("r", r))
			return r, nil
		}
	}

	return nil, nil
}
