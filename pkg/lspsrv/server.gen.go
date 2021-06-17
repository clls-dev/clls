// Do not edit, generated by github.com/clls-dev/clls/pkg/lspsrv/gen

package lspsrv

import (
	"encoding/json"

	"github.com/clls-dev/clls/pkg/lsp"
	"github.com/pkg/errors"
)

type LanguageServer interface {
	Initialized(*lsp.InitializedParams) error
	Exit() error
	DidOpenTextDocument(*lsp.DidOpenTextDocumentParams) error
	DidCloseTextDocument(*lsp.DidCloseTextDocumentParams) error
	DidChangeTextDocument(*lsp.DidChangeTextDocumentParams) error
	DidSaveTextDocument(*lsp.DidSaveTextDocumentParams) error
	Initialize(*lsp.InitializeParams) (*lsp.InitializeResult, error)
	Shutdown() error
	SemanticTokens(*lsp.SemanticTokensParams) (*lsp.SemanticTokens, error)
	DocumentFormatting(*lsp.DocumentFormattingParams) ([]lsp.TextEdit, error)
	Rename(*lsp.RenameParams) (*lsp.WorkspaceEdit, error)
}

type UnimplementedLanguageServer struct{}

var _ LanguageServer = (*UnimplementedLanguageServer)(nil)

func (s *UnimplementedLanguageServer) Initialized(*lsp.InitializedParams) error {
	return ErrNotImplemented
}

func (s *UnimplementedLanguageServer) Exit() error {
	return ErrNotImplemented
}

func (s *UnimplementedLanguageServer) DidOpenTextDocument(*lsp.DidOpenTextDocumentParams) error {
	return ErrNotImplemented
}

func (s *UnimplementedLanguageServer) DidCloseTextDocument(*lsp.DidCloseTextDocumentParams) error {
	return ErrNotImplemented
}

func (s *UnimplementedLanguageServer) DidChangeTextDocument(*lsp.DidChangeTextDocumentParams) error {
	return ErrNotImplemented
}

func (s *UnimplementedLanguageServer) DidSaveTextDocument(*lsp.DidSaveTextDocumentParams) error {
	return ErrNotImplemented
}

func (s *UnimplementedLanguageServer) Initialize(*lsp.InitializeParams) (*lsp.InitializeResult, error) {
	return nil, ErrNotImplemented
}

func (s *UnimplementedLanguageServer) Shutdown() error {
	return ErrNotImplemented
}

func (s *UnimplementedLanguageServer) SemanticTokens(*lsp.SemanticTokensParams) (*lsp.SemanticTokens, error) {
	return nil, ErrNotImplemented
}

func (s *UnimplementedLanguageServer) DocumentFormatting(*lsp.DocumentFormattingParams) ([]lsp.TextEdit, error) {
	return nil, ErrNotImplemented
}

func (s *UnimplementedLanguageServer) Rename(*lsp.RenameParams) (*lsp.WorkspaceEdit, error) {
	return nil, ErrNotImplemented
}

func LanguageServerHandle(s LanguageServer, method string, payloadBytes []byte) (interface{}, error) {
	switch method {
	case "initialized":
		var payload lsp.InitializedParams
		if err := json.Unmarshal(payloadBytes, &payload); err != nil {
			return nil, errors.Wrap(err, "unmarshal payload")
		}
		return nil, s.Initialized(&payload)

	case "exit":
		return nil, s.Exit()

	case "textDocument/didOpen":
		var payload lsp.DidOpenTextDocumentParams
		if err := json.Unmarshal(payloadBytes, &payload); err != nil {
			return nil, errors.Wrap(err, "unmarshal payload")
		}
		return nil, s.DidOpenTextDocument(&payload)

	case "textDocument/didClose":
		var payload lsp.DidCloseTextDocumentParams
		if err := json.Unmarshal(payloadBytes, &payload); err != nil {
			return nil, errors.Wrap(err, "unmarshal payload")
		}
		return nil, s.DidCloseTextDocument(&payload)

	case "textDocument/didChange":
		var payload lsp.DidChangeTextDocumentParams
		if err := json.Unmarshal(payloadBytes, &payload); err != nil {
			return nil, errors.Wrap(err, "unmarshal payload")
		}
		return nil, s.DidChangeTextDocument(&payload)

	case "textDocument/didSave":
		var payload lsp.DidSaveTextDocumentParams
		if err := json.Unmarshal(payloadBytes, &payload); err != nil {
			return nil, errors.Wrap(err, "unmarshal payload")
		}
		return nil, s.DidSaveTextDocument(&payload)

	case "initialize":
		var payload lsp.InitializeParams
		if err := json.Unmarshal(payloadBytes, &payload); err != nil {
			return nil, errors.Wrap(err, "unmarshal payload")
		}
		return s.Initialize(&payload)

	case "shutdown":
		return nil, s.Shutdown()

	case "textDocument/semanticTokens/full":
		var payload lsp.SemanticTokensParams
		if err := json.Unmarshal(payloadBytes, &payload); err != nil {
			return nil, errors.Wrap(err, "unmarshal payload")
		}
		return s.SemanticTokens(&payload)

	case "textDocument/formatting":
		var payload lsp.DocumentFormattingParams
		if err := json.Unmarshal(payloadBytes, &payload); err != nil {
			return nil, errors.Wrap(err, "unmarshal payload")
		}
		return s.DocumentFormatting(&payload)

	case "textDocument/rename":
		var payload lsp.RenameParams
		if err := json.Unmarshal(payloadBytes, &payload); err != nil {
			return nil, errors.Wrap(err, "unmarshal payload")
		}
		return s.Rename(&payload)

	}
	return nil, ErrUnknownMethod
}
