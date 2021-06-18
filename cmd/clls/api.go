package main

import (
	"encoding/base64"

	"github.com/clls-dev/clls/pkg/clls"
	"github.com/clls-dev/clls/pkg/lsp"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/crypto/sha3"
)

func (s *server) Initialize(*lsp.InitializeParams) (*lsp.InitializeResult, error) {
	caps := lsp.ServerCapabilities{
		TextDocumentSync: lsp.Full,
		SemanticTokensProvider: lsp.SemanticTokensOptions{
			Legend: lsp.StandardSemanticTokensLegend,
			Full:   true,
		},
		DocumentFormattingProvider: true,
		RenameProvider:             true,
		DocumentHighlightProvider:  true,
	}
	s.l.Debug("server initialized", zap.Any("capabilities", caps))
	return &lsp.InitializeResult{
		ServerInfo: &lsp.InitializeResultServerInfo{
			Name:    "clls",
			Version: "0.1.0",
		},
		Capabilities: &caps}, nil

}

func hashString(s string) string {
	h := sha3.New256()
	if _, err := h.Write([]byte(s)); err != nil {
		panic(err)
	}
	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}

func newDocumentData(text string) *documentData {
	return &documentData{
		content:     text,
		contentHash: hashString(text), // FIXME, if an include changes, the cache will be bad
	}
}

func (s *server) DidOpenTextDocument(params *lsp.DidOpenTextDocumentParams) error {
	docData := newDocumentData(params.TextDocument.Text)
	if pulled, ok := s.cache.pull(docData.contentHash); ok {
		s.openedDocs[params.TextDocument.URI] = pulled
		return nil
	}
	s.openedDocs[params.TextDocument.URI] = newDocumentData(params.TextDocument.Text)
	return nil
}

// This only supports full file changes
func (s *server) DidChangeTextDocument(params *lsp.DidChangeTextDocumentParams) error {
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

func (s *server) DidCloseTextDocument(params *lsp.DidCloseTextDocumentParams) error {
	dd, ok := s.openedDocs[params.TextDocument.URI]
	if !ok {
		return nil
	}
	delete(s.openedDocs, params.TextDocument.URI)
	s.cache.put(dd)
	return nil
}

func (s *server) Rename(params *lsp.RenameParams) (*lsp.WorkspaceEdit, error) {
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
				Changes: map[string][]lsp.TextEdit{},
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

func (s *server) DocumentFormatting(params *lsp.DocumentFormattingParams) ([]lsp.TextEdit, error) {
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
			End:   lsp.Position{Line: lsp.UInteger(linesCount + 1), Character: 0},
		},
		NewText: newText,
	}}, nil
}

func (s *server) SemanticTokens(params *lsp.SemanticTokensParams) (*lsp.SemanticTokens, error) {
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

func (s *server) Shutdown() error {
	s.down = true
	return nil
}

func (s *server) Exit() error {
	s.exit = true
	return nil
}

func (s *server) Initialized(*lsp.InitializedParams) error {
	return nil
}

func (s *server) DidSaveTextDocument(*lsp.DidSaveTextDocumentParams) error {
	return nil // this is to gracefully ignore the event
}

func (s *server) DocumentHighlight(params *lsp.DocumentHighlightParams) ([]lsp.DocumentHighlight, error) {
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
				r = append(r, lsp.DocumentHighlight{Range: t.Range(), Kind: &lsp.Text})
			}
			s.l.Debug("will send", zap.Any("r", r))
			return r, nil
		}
	}

	return nil, errors.New("symbol not found")
}
