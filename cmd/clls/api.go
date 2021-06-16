package main

import (
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"

	"github.com/clls-dev/clls/pkg/clls"
	"github.com/clls-dev/clls/pkg/lsp"
	"github.com/clls-dev/clls/pkg/lsph"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func (s *server) Initialize(*lsp.InitializeParams) (*lsp.InitializeResult, error) {
	return &lsp.InitializeResult{
		ServerInfo: &lsp.InitializeResultServerInfo{
			Name:    "clls",
			Version: "0.1.0",
		},
		Capabilities: &lsp.ServerCapabilities{
			TextDocumentSync: lsp.Full,
			SemanticTokensProvider: lsp.SemanticTokensOptions{
				Legend: lsp.StandardSemanticTokensLegend,
				Full:   true,
			},
			DocumentFormattingProvider: true,
			RenameProvider:             true,
		},
	}, nil
}

func (s *server) DidOpenTextDocument(params *lsp.DidOpenTextDocumentParams) error {
	s.docs[params.TextDocument.URI] = params.TextDocument.Text
	return nil
}

func (s *server) DidChangeTextDocument(params *lsp.DidChangeTextDocumentParams) error {
	if len(params.ContentChanges) == 0 {
		return nil
	}
	s.docs[params.TextDocument.URI] = params.ContentChanges[0].Text
	return nil
}

func (s *server) DidCloseTextDocument(params *lsp.DidCloseTextDocumentParams) error {
	delete(s.docs, params.TextDocument.URI)
	return nil
}

func (s *server) Rename(params *lsp.RenameParams) (*lsp.WorkspaceEdit, error) {
	uriStr := params.TextDocument.URI
	s.l.Debug("textDocument/rename " + uriStr)
	if !strings.HasPrefix(uriStr, fileURIPrefix) {
		return nil, errors.New("textDocument.uri is not a file uri")
	}
	pathStr := strings.TrimPrefix(uriStr, fileURIPrefix)
	dirname := filepath.Dir(pathStr)
	filename := filepath.Base(pathStr)

	mod, err := s.loadCLVM(filename, dirname, uriStr)
	if err != nil {
		return nil, errors.Wrap(err, "parse module")
	}

	p := params.Position
	line := int(p.Line)
	char := int(p.Character)
	newName := params.NewName

	ss := mod.Symbols(s.l)
	s.l.Debug("got symbols", zap.Any("ss", ss))

	for _, s := range ss {
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
	s.l.Debug("textDocument/formatting " + uriStr)
	if !strings.HasPrefix(uriStr, fileURIPrefix) {
		return nil, errors.New("textDocument.uri is not a file uri")
	}
	pathStr := strings.TrimPrefix(uriStr, fileURIPrefix)
	s.l.Debug("parsed uri", zap.String("path", pathStr))
	fileStr := ""
	if s, ok := s.docs[uriStr]; ok {
		fileStr = s
	} else {
		b, err := ioutil.ReadFile(pathStr)
		if err != nil {
			return nil, errors.Wrap(err, "read file")
		}
		fileStr = string(b)
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
	uriStr := params.TextDocument.URI
	s.l.Debug("textDocument/formatting " + uriStr)
	if !strings.HasPrefix(uriStr, fileURIPrefix) {
		return nil, errors.New("textDocument.uri is not a file uri")
	}
	pathStr := strings.TrimPrefix(uriStr, fileURIPrefix)
	s.l.Debug("parsed uri", zap.String("path", pathStr))
	dirname := filepath.Dir(pathStr)
	filename := filepath.Base(pathStr)
	mod, err := s.loadCLVM(filename, dirname, uriStr)
	if err != nil {
		return nil, errors.Wrap(err, "parse module")
	}
	//l.Debug("parsed mod", zap.Any("val", mod))

	inserts := []insert(nil)

	if mod.IsMod {
		inserts = append(inserts, insert{Kind: "keyword", Token: mod.ModToken})
		if mod.Args != nil {
			inserts = insertParamsTokens(inserts, mod.Args)
		}
	}

	if mod.Constants != nil {
		for _, c := range mod.Constants {
			if c.Token != nil {
				inserts = append(inserts, insert{Kind: "keyword", Token: c.Token})
			}
			if t, ok := c.Name.(*clls.Token); ok && t != nil {
				inserts = append(inserts, insert{Kind: "variable", Modifiers: []string{"readonly"}, Token: t})
			}
			inserts = insertBody(inserts, c.Value, clls.BuiltinFuncsByName)
		}
	}

	for _, t := range mod.Comments {
		inserts = append(inserts, insert{Kind: "comment", Token: t})
	}

	for _, incl := range mod.Includes {
		inserts = append(inserts, insert{Kind: "keyword", Token: incl.Token})
		if t, ok := incl.Value.(*clls.Token); ok && t != nil {
			inserts = append(inserts, insert{Kind: "string", Token: t})
		}
	}

	allFuncs := map[string]*clls.Function{}
	for k, v := range clls.BuiltinFuncsByName {
		allFuncs[k] = v
	}
	for k, v := range mod.FunctionsByName {
		allFuncs[k] = v
	}
	if len(mod.Functions) > 0 {
		for _, f := range mod.Functions {
			inserts = append(inserts, insert{Kind: "keyword", Token: f.KeywordToken})
			if f.Name != nil {
				inserts = append(inserts, insert{Kind: "function", Token: f.Name})
			}
			if f.Params != nil {
				inserts = insertParamsTokens(inserts, f.Params)
			}
			inserts = insertBody(inserts, f.Body, allFuncs)
		}
	}

	if mod.Main != nil {
		inserts = insertBody(inserts, mod.Main, allFuncs)
	}

	data := []lsp.UInteger(nil)

	if len(inserts) != 0 {
		ninserts := inserts
		inserts := []insert(nil)
		for _, i := range ninserts {
			if i.Token != nil {
				inserts = append(inserts, i)
			}
		}

		sort.Slice(inserts, func(i, j int) bool {
			a := inserts[i]
			b := inserts[j]
			ai := 0
			if a.Token != nil {
				ai = a.Token.Index
			}
			bi := 0
			if b.Token != nil {
				bi = b.Token.Index
			}
			return ai < bi
		})

		ltoks := lsph.SemanticTokenSlice{}

		t := inserts[0].Token
		if t != nil {
			tt, tm, err := tokenInfo(s.l, inserts[0].Kind, inserts[0].Modifiers, &lsp.StandardSemanticTokensLegend)
			if err != nil {
				panic(err)
			}
			ltoks = append(ltoks, lsph.SemanticToken{
				DeltaLine:      lsp.UInteger(t.Line),
				DeltaStartChar: lsp.UInteger(t.StartChar),
				Length:         lsp.UInteger(len(t.Text)),
				TokenType:      tt,
				TokenModifiers: tm,
			})

			for i := 1; i < len(inserts); i++ {
				prev := inserts[i-1]
				in := inserts[i]
				t := in.Token
				pt := prev.Token
				deltaLine := t.Line - pt.Line
				if deltaLine < 0 {
					panic("negative line delta")
				}
				deltaStartChar := t.StartChar
				if deltaLine == 0 {
					deltaStartChar = t.StartChar - pt.StartChar
				}
				if deltaLine < 0 {
					panic("negative start char delta")
				}

				tt, tm, err := tokenInfo(s.l, in.Kind, in.Modifiers, &lsp.StandardSemanticTokensLegend)
				if err != nil {
					panic(err)
				}
				ltoks = append(ltoks, lsph.SemanticToken{
					DeltaLine:      lsp.UInteger(deltaLine),
					DeltaStartChar: lsp.UInteger(deltaStartChar),
					Length:         lsp.UInteger(len(t.Text)),
					TokenType:      tt,
					TokenModifiers: tm,
				})
			}

			data = ltoks.Flat()
		}
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
