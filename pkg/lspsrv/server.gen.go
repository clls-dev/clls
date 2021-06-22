package lspsrv

import (
	"context"
	"encoding/json"

	lsp "go.lsp.dev/protocol"
)

type UnimplementedLanguageServer struct{}

var _ lsp.Server = (*UnimplementedLanguageServer)(nil)

func (*UnimplementedLanguageServer) CodeAction(context.Context, *lsp.CodeActionParams) ([]lsp.CodeAction, error) {
	return nil, ErrNotImplemented
}

func (*UnimplementedLanguageServer) CodeLens(context.Context, *lsp.CodeLensParams) ([]lsp.CodeLens, error) {
	return nil, ErrNotImplemented
}

func (*UnimplementedLanguageServer) CodeLensRefresh(context.Context) error {
	return ErrNotImplemented
}

func (*UnimplementedLanguageServer) CodeLensResolve(context.Context, *lsp.CodeLens) (*lsp.CodeLens, error) {
	return nil, ErrNotImplemented
}

func (*UnimplementedLanguageServer) ColorPresentation(context.Context, *lsp.ColorPresentationParams) ([]lsp.ColorPresentation, error) {
	return nil, ErrNotImplemented
}

func (*UnimplementedLanguageServer) Completion(context.Context, *lsp.CompletionParams) (*lsp.CompletionList, error) {
	return nil, ErrNotImplemented
}

func (*UnimplementedLanguageServer) CompletionResolve(context.Context, *lsp.CompletionItem) (*lsp.CompletionItem, error) {
	return nil, ErrNotImplemented
}

func (*UnimplementedLanguageServer) Declaration(context.Context, *lsp.DeclarationParams) ([]lsp.Location, error) {
	return nil, ErrNotImplemented
}

func (*UnimplementedLanguageServer) Definition(context.Context, *lsp.DefinitionParams) ([]lsp.Location, error) {
	return nil, ErrNotImplemented
}

func (*UnimplementedLanguageServer) DidChange(context.Context, *lsp.DidChangeTextDocumentParams) error {
	return ErrNotImplemented
}

func (*UnimplementedLanguageServer) DidChangeConfiguration(context.Context, *lsp.DidChangeConfigurationParams) error {
	return ErrNotImplemented
}

func (*UnimplementedLanguageServer) DidChangeWatchedFiles(context.Context, *lsp.DidChangeWatchedFilesParams) error {
	return ErrNotImplemented
}

func (*UnimplementedLanguageServer) DidChangeWorkspaceFolders(context.Context, *lsp.DidChangeWorkspaceFoldersParams) error {
	return ErrNotImplemented
}

func (*UnimplementedLanguageServer) DidClose(context.Context, *lsp.DidCloseTextDocumentParams) error {
	return ErrNotImplemented
}

func (*UnimplementedLanguageServer) DidCreateFiles(context.Context, *lsp.CreateFilesParams) error {
	return ErrNotImplemented
}

func (*UnimplementedLanguageServer) DidDeleteFiles(context.Context, *lsp.DeleteFilesParams) error {
	return ErrNotImplemented
}

func (*UnimplementedLanguageServer) DidOpen(context.Context, *lsp.DidOpenTextDocumentParams) error {
	return ErrNotImplemented
}

func (*UnimplementedLanguageServer) DidRenameFiles(context.Context, *lsp.RenameFilesParams) error {
	return ErrNotImplemented
}

func (*UnimplementedLanguageServer) DidSave(context.Context, *lsp.DidSaveTextDocumentParams) error {
	return ErrNotImplemented
}

func (*UnimplementedLanguageServer) DocumentColor(context.Context, *lsp.DocumentColorParams) ([]lsp.ColorInformation, error) {
	return nil, ErrNotImplemented
}

func (*UnimplementedLanguageServer) DocumentHighlight(context.Context, *lsp.DocumentHighlightParams) ([]lsp.DocumentHighlight, error) {
	return nil, ErrNotImplemented
}

func (*UnimplementedLanguageServer) DocumentLink(context.Context, *lsp.DocumentLinkParams) ([]lsp.DocumentLink, error) {
	return nil, ErrNotImplemented
}

func (*UnimplementedLanguageServer) DocumentLinkResolve(context.Context, *lsp.DocumentLink) (*lsp.DocumentLink, error) {
	return nil, ErrNotImplemented
}

func (*UnimplementedLanguageServer) DocumentSymbol(context.Context, *lsp.DocumentSymbolParams) ([]interface{}, error) {
	return nil, ErrNotImplemented
}

func (*UnimplementedLanguageServer) ExecuteCommand(context.Context, *lsp.ExecuteCommandParams) (interface{}, error) {
	return nil, ErrNotImplemented
}

func (*UnimplementedLanguageServer) Exit(context.Context) error {
	return ErrNotImplemented
}

func (*UnimplementedLanguageServer) FoldingRanges(context.Context, *lsp.FoldingRangeParams) ([]lsp.FoldingRange, error) {
	return nil, ErrNotImplemented
}

func (*UnimplementedLanguageServer) Formatting(context.Context, *lsp.DocumentFormattingParams) ([]lsp.TextEdit, error) {
	return nil, ErrNotImplemented
}

func (*UnimplementedLanguageServer) Hover(context.Context, *lsp.HoverParams) (*lsp.Hover, error) {
	return nil, ErrNotImplemented
}

func (*UnimplementedLanguageServer) Implementation(context.Context, *lsp.ImplementationParams) ([]lsp.Location, error) {
	return nil, ErrNotImplemented
}

func (*UnimplementedLanguageServer) IncomingCalls(context.Context, *lsp.CallHierarchyIncomingCallsParams) ([]lsp.CallHierarchyIncomingCall, error) {
	return nil, ErrNotImplemented
}

func (*UnimplementedLanguageServer) Initialize(context.Context, *lsp.InitializeParams) (*lsp.InitializeResult, error) {
	return nil, ErrNotImplemented
}

func (*UnimplementedLanguageServer) Initialized(context.Context, *lsp.InitializedParams) error {
	return ErrNotImplemented
}

func (*UnimplementedLanguageServer) LinkedEditingRange(context.Context, *lsp.LinkedEditingRangeParams) (*lsp.LinkedEditingRanges, error) {
	return nil, ErrNotImplemented
}

func (*UnimplementedLanguageServer) LogTrace(context.Context, *lsp.LogTraceParams) error {
	return ErrNotImplemented
}

func (*UnimplementedLanguageServer) Moniker(context.Context, *lsp.MonikerParams) ([]lsp.Moniker, error) {
	return nil, ErrNotImplemented
}

func (*UnimplementedLanguageServer) OnTypeFormatting(context.Context, *lsp.DocumentOnTypeFormattingParams) ([]lsp.TextEdit, error) {
	return nil, ErrNotImplemented
}

func (*UnimplementedLanguageServer) OutgoingCalls(context.Context, *lsp.CallHierarchyOutgoingCallsParams) ([]lsp.CallHierarchyOutgoingCall, error) {
	return nil, ErrNotImplemented
}

func (*UnimplementedLanguageServer) PrepareCallHierarchy(context.Context, *lsp.CallHierarchyPrepareParams) ([]lsp.CallHierarchyItem, error) {
	return nil, ErrNotImplemented
}

func (*UnimplementedLanguageServer) PrepareRename(context.Context, *lsp.PrepareRenameParams) (*lsp.Range, error) {
	return nil, ErrNotImplemented
}

func (*UnimplementedLanguageServer) RangeFormatting(context.Context, *lsp.DocumentRangeFormattingParams) ([]lsp.TextEdit, error) {
	return nil, ErrNotImplemented
}

func (*UnimplementedLanguageServer) References(context.Context, *lsp.ReferenceParams) ([]lsp.Location, error) {
	return nil, ErrNotImplemented
}

func (*UnimplementedLanguageServer) Rename(context.Context, *lsp.RenameParams) (*lsp.WorkspaceEdit, error) {
	return nil, ErrNotImplemented
}

func (*UnimplementedLanguageServer) Request(context.Context, string, interface{}) (interface{}, error) {
	return nil, ErrNotImplemented
}

func (*UnimplementedLanguageServer) SemanticTokensFull(context.Context, *lsp.SemanticTokensParams) (*lsp.SemanticTokens, error) {
	return nil, ErrNotImplemented
}

func (*UnimplementedLanguageServer) SemanticTokensFullDelta(context.Context, *lsp.SemanticTokensDeltaParams) (interface{}, error) {
	return nil, ErrNotImplemented
}

func (*UnimplementedLanguageServer) SemanticTokensRange(context.Context, *lsp.SemanticTokensRangeParams) (*lsp.SemanticTokens, error) {
	return nil, ErrNotImplemented
}

func (*UnimplementedLanguageServer) SemanticTokensRefresh(context.Context) error {
	return ErrNotImplemented
}

func (*UnimplementedLanguageServer) SetTrace(context.Context, *lsp.SetTraceParams) error {
	return ErrNotImplemented
}

func (*UnimplementedLanguageServer) ShowDocument(context.Context, *lsp.ShowDocumentParams) (*lsp.ShowDocumentResult, error) {
	return nil, ErrNotImplemented
}

func (*UnimplementedLanguageServer) Shutdown(context.Context) error {
	return ErrNotImplemented
}

func (*UnimplementedLanguageServer) SignatureHelp(context.Context, *lsp.SignatureHelpParams) (*lsp.SignatureHelp, error) {
	return nil, ErrNotImplemented
}

func (*UnimplementedLanguageServer) Symbols(context.Context, *lsp.WorkspaceSymbolParams) ([]lsp.SymbolInformation, error) {
	return nil, ErrNotImplemented
}

func (*UnimplementedLanguageServer) TypeDefinition(context.Context, *lsp.TypeDefinitionParams) ([]lsp.Location, error) {
	return nil, ErrNotImplemented
}

func (*UnimplementedLanguageServer) WillCreateFiles(context.Context, *lsp.CreateFilesParams) (*lsp.WorkspaceEdit, error) {
	return nil, ErrNotImplemented
}

func (*UnimplementedLanguageServer) WillDeleteFiles(context.Context, *lsp.DeleteFilesParams) (*lsp.WorkspaceEdit, error) {
	return nil, ErrNotImplemented
}

func (*UnimplementedLanguageServer) WillRenameFiles(context.Context, *lsp.RenameFilesParams) (*lsp.WorkspaceEdit, error) {
	return nil, ErrNotImplemented
}

func (*UnimplementedLanguageServer) WillSave(context.Context, *lsp.WillSaveTextDocumentParams) error {
	return ErrNotImplemented
}

func (*UnimplementedLanguageServer) WillSaveWaitUntil(context.Context, *lsp.WillSaveTextDocumentParams) ([]lsp.TextEdit, error) {
	return nil, ErrNotImplemented
}

func (*UnimplementedLanguageServer) WorkDoneProgressCancel(context.Context, *lsp.WorkDoneProgressCancelParams) error {
	return ErrNotImplemented
}
func Unmarshal(method string, payloadBytes []byte) (interface{}, error) {
	switch method {
	case "textDocument/codeAction":
		var payload lsp.CodeActionParams
		return &payload, json.Unmarshal(payloadBytes, &payload)

	case "textDocument/codeLens":
		var payload lsp.CodeLensParams
		return &payload, json.Unmarshal(payloadBytes, &payload)

	case "textDocument/codeLensResolve":
		var payload lsp.CodeLens
		return &payload, json.Unmarshal(payloadBytes, &payload)

	case "textDocument/colorPresentation":
		var payload lsp.ColorPresentationParams
		return &payload, json.Unmarshal(payloadBytes, &payload)

	case "textDocument/completion":
		var payload lsp.CompletionParams
		return &payload, json.Unmarshal(payloadBytes, &payload)

	case "textDocument/completionResolve":
		var payload lsp.CompletionItem
		return &payload, json.Unmarshal(payloadBytes, &payload)

	case "textDocument/declaration":
		var payload lsp.DeclarationParams
		return &payload, json.Unmarshal(payloadBytes, &payload)

	case "textDocument/definition":
		var payload lsp.DefinitionParams
		return &payload, json.Unmarshal(payloadBytes, &payload)

	case "textDocument/didChange":
		var payload lsp.DidChangeTextDocumentParams
		return &payload, json.Unmarshal(payloadBytes, &payload)

	case "textDocument/didChangeConfiguration":
		var payload lsp.DidChangeConfigurationParams
		return &payload, json.Unmarshal(payloadBytes, &payload)

	case "textDocument/didChangeWatchedFiles":
		var payload lsp.DidChangeWatchedFilesParams
		return &payload, json.Unmarshal(payloadBytes, &payload)

	case "textDocument/didChangeWorkspaceFolders":
		var payload lsp.DidChangeWorkspaceFoldersParams
		return &payload, json.Unmarshal(payloadBytes, &payload)

	case "textDocument/didClose":
		var payload lsp.DidCloseTextDocumentParams
		return &payload, json.Unmarshal(payloadBytes, &payload)

	case "textDocument/didCreateFiles":
		var payload lsp.CreateFilesParams
		return &payload, json.Unmarshal(payloadBytes, &payload)

	case "textDocument/didDeleteFiles":
		var payload lsp.DeleteFilesParams
		return &payload, json.Unmarshal(payloadBytes, &payload)

	case "textDocument/didOpen":
		var payload lsp.DidOpenTextDocumentParams
		return &payload, json.Unmarshal(payloadBytes, &payload)

	case "textDocument/didRenameFiles":
		var payload lsp.RenameFilesParams
		return &payload, json.Unmarshal(payloadBytes, &payload)

	case "textDocument/didSave":
		var payload lsp.DidSaveTextDocumentParams
		return &payload, json.Unmarshal(payloadBytes, &payload)

	case "textDocument/documentColor":
		var payload lsp.DocumentColorParams
		return &payload, json.Unmarshal(payloadBytes, &payload)

	case "textDocument/documentHighlight":
		var payload lsp.DocumentHighlightParams
		return &payload, json.Unmarshal(payloadBytes, &payload)

	case "textDocument/documentLink":
		var payload lsp.DocumentLinkParams
		return &payload, json.Unmarshal(payloadBytes, &payload)

	case "textDocument/documentLinkResolve":
		var payload lsp.DocumentLink
		return &payload, json.Unmarshal(payloadBytes, &payload)

	case "textDocument/documentSymbol":
		var payload lsp.DocumentSymbolParams
		return &payload, json.Unmarshal(payloadBytes, &payload)

	case "textDocument/executeCommand":
		var payload lsp.ExecuteCommandParams
		return &payload, json.Unmarshal(payloadBytes, &payload)

	case "textDocument/foldingRanges":
		var payload lsp.FoldingRangeParams
		return &payload, json.Unmarshal(payloadBytes, &payload)

	case "textDocument/formatting":
		var payload lsp.DocumentFormattingParams
		return &payload, json.Unmarshal(payloadBytes, &payload)

	case "textDocument/hover":
		var payload lsp.HoverParams
		return &payload, json.Unmarshal(payloadBytes, &payload)

	case "textDocument/implementation":
		var payload lsp.ImplementationParams
		return &payload, json.Unmarshal(payloadBytes, &payload)

	case "textDocument/incomingCalls":
		var payload lsp.CallHierarchyIncomingCallsParams
		return &payload, json.Unmarshal(payloadBytes, &payload)

	case "initialize":
		var payload lsp.InitializeParams
		return &payload, json.Unmarshal(payloadBytes, &payload)

	case "initialized":
		var payload lsp.InitializedParams
		return &payload, json.Unmarshal(payloadBytes, &payload)

	case "textDocument/linkedEditingRange":
		var payload lsp.LinkedEditingRangeParams
		return &payload, json.Unmarshal(payloadBytes, &payload)

	case "textDocument/logTrace":
		var payload lsp.LogTraceParams
		return &payload, json.Unmarshal(payloadBytes, &payload)

	case "textDocument/moniker":
		var payload lsp.MonikerParams
		return &payload, json.Unmarshal(payloadBytes, &payload)

	case "textDocument/onTypeFormatting":
		var payload lsp.DocumentOnTypeFormattingParams
		return &payload, json.Unmarshal(payloadBytes, &payload)

	case "textDocument/outgoingCalls":
		var payload lsp.CallHierarchyOutgoingCallsParams
		return &payload, json.Unmarshal(payloadBytes, &payload)

	case "textDocument/prepareCallHierarchy":
		var payload lsp.CallHierarchyPrepareParams
		return &payload, json.Unmarshal(payloadBytes, &payload)

	case "textDocument/prepareRename":
		var payload lsp.PrepareRenameParams
		return &payload, json.Unmarshal(payloadBytes, &payload)

	case "textDocument/rangeFormatting":
		var payload lsp.DocumentRangeFormattingParams
		return &payload, json.Unmarshal(payloadBytes, &payload)

	case "textDocument/references":
		var payload lsp.ReferenceParams
		return &payload, json.Unmarshal(payloadBytes, &payload)

	case "textDocument/rename":
		var payload lsp.RenameParams
		return &payload, json.Unmarshal(payloadBytes, &payload)

	case "textDocument/semanticTokens/full":
		var payload lsp.SemanticTokensParams
		return &payload, json.Unmarshal(payloadBytes, &payload)

	case "textDocument/semanticTokensFullDelta":
		var payload lsp.SemanticTokensDeltaParams
		return &payload, json.Unmarshal(payloadBytes, &payload)

	case "textDocument/semanticTokensRange":
		var payload lsp.SemanticTokensRangeParams
		return &payload, json.Unmarshal(payloadBytes, &payload)

	case "textDocument/setTrace":
		var payload lsp.SetTraceParams
		return &payload, json.Unmarshal(payloadBytes, &payload)

	case "textDocument/showDocument":
		var payload lsp.ShowDocumentParams
		return &payload, json.Unmarshal(payloadBytes, &payload)

	case "textDocument/signatureHelp":
		var payload lsp.SignatureHelpParams
		return &payload, json.Unmarshal(payloadBytes, &payload)

	case "textDocument/symbols":
		var payload lsp.WorkspaceSymbolParams
		return &payload, json.Unmarshal(payloadBytes, &payload)

	case "textDocument/typeDefinition":
		var payload lsp.TypeDefinitionParams
		return &payload, json.Unmarshal(payloadBytes, &payload)

	case "textDocument/willCreateFiles":
		var payload lsp.CreateFilesParams
		return &payload, json.Unmarshal(payloadBytes, &payload)

	case "textDocument/willDeleteFiles":
		var payload lsp.DeleteFilesParams
		return &payload, json.Unmarshal(payloadBytes, &payload)

	case "textDocument/willRenameFiles":
		var payload lsp.RenameFilesParams
		return &payload, json.Unmarshal(payloadBytes, &payload)

	case "textDocument/willSave":
		var payload lsp.WillSaveTextDocumentParams
		return &payload, json.Unmarshal(payloadBytes, &payload)

	case "textDocument/willSaveWaitUntil":
		var payload lsp.WillSaveTextDocumentParams
		return &payload, json.Unmarshal(payloadBytes, &payload)

	case "textDocument/workDoneProgressCancel":
		var payload lsp.WorkDoneProgressCancelParams
		return &payload, json.Unmarshal(payloadBytes, &payload)

	}
	return nil, ErrUnknownMethod
}

func Request(ctx context.Context, s lsp.Server, method string, payload interface{}) (interface{}, error) {
	switch method {
	case "textDocument/codeAction":
		castedPayload, ok := payload.(*lsp.CodeActionParams)
		if !ok {
			return nil, ErrBadPayloadType
		}
		return s.CodeAction(ctx, castedPayload)

	case "textDocument/codeLens":
		castedPayload, ok := payload.(*lsp.CodeLensParams)
		if !ok {
			return nil, ErrBadPayloadType
		}
		return s.CodeLens(ctx, castedPayload)

	case "textDocument/codeLensRefresh":
		return nil, s.CodeLensRefresh(ctx)

	case "textDocument/codeLensResolve":
		castedPayload, ok := payload.(*lsp.CodeLens)
		if !ok {
			return nil, ErrBadPayloadType
		}
		return s.CodeLensResolve(ctx, castedPayload)

	case "textDocument/colorPresentation":
		castedPayload, ok := payload.(*lsp.ColorPresentationParams)
		if !ok {
			return nil, ErrBadPayloadType
		}
		return s.ColorPresentation(ctx, castedPayload)

	case "textDocument/completion":
		castedPayload, ok := payload.(*lsp.CompletionParams)
		if !ok {
			return nil, ErrBadPayloadType
		}
		return s.Completion(ctx, castedPayload)

	case "textDocument/completionResolve":
		castedPayload, ok := payload.(*lsp.CompletionItem)
		if !ok {
			return nil, ErrBadPayloadType
		}
		return s.CompletionResolve(ctx, castedPayload)

	case "textDocument/declaration":
		castedPayload, ok := payload.(*lsp.DeclarationParams)
		if !ok {
			return nil, ErrBadPayloadType
		}
		return s.Declaration(ctx, castedPayload)

	case "textDocument/definition":
		castedPayload, ok := payload.(*lsp.DefinitionParams)
		if !ok {
			return nil, ErrBadPayloadType
		}
		return s.Definition(ctx, castedPayload)

	case "textDocument/didChange":
		castedPayload, ok := payload.(*lsp.DidChangeTextDocumentParams)
		if !ok {
			return nil, ErrBadPayloadType
		}
		return nil, s.DidChange(ctx, castedPayload)

	case "textDocument/didChangeConfiguration":
		castedPayload, ok := payload.(*lsp.DidChangeConfigurationParams)
		if !ok {
			return nil, ErrBadPayloadType
		}
		return nil, s.DidChangeConfiguration(ctx, castedPayload)

	case "textDocument/didChangeWatchedFiles":
		castedPayload, ok := payload.(*lsp.DidChangeWatchedFilesParams)
		if !ok {
			return nil, ErrBadPayloadType
		}
		return nil, s.DidChangeWatchedFiles(ctx, castedPayload)

	case "textDocument/didChangeWorkspaceFolders":
		castedPayload, ok := payload.(*lsp.DidChangeWorkspaceFoldersParams)
		if !ok {
			return nil, ErrBadPayloadType
		}
		return nil, s.DidChangeWorkspaceFolders(ctx, castedPayload)

	case "textDocument/didClose":
		castedPayload, ok := payload.(*lsp.DidCloseTextDocumentParams)
		if !ok {
			return nil, ErrBadPayloadType
		}
		return nil, s.DidClose(ctx, castedPayload)

	case "textDocument/didCreateFiles":
		castedPayload, ok := payload.(*lsp.CreateFilesParams)
		if !ok {
			return nil, ErrBadPayloadType
		}
		return nil, s.DidCreateFiles(ctx, castedPayload)

	case "textDocument/didDeleteFiles":
		castedPayload, ok := payload.(*lsp.DeleteFilesParams)
		if !ok {
			return nil, ErrBadPayloadType
		}
		return nil, s.DidDeleteFiles(ctx, castedPayload)

	case "textDocument/didOpen":
		castedPayload, ok := payload.(*lsp.DidOpenTextDocumentParams)
		if !ok {
			return nil, ErrBadPayloadType
		}
		return nil, s.DidOpen(ctx, castedPayload)

	case "textDocument/didRenameFiles":
		castedPayload, ok := payload.(*lsp.RenameFilesParams)
		if !ok {
			return nil, ErrBadPayloadType
		}
		return nil, s.DidRenameFiles(ctx, castedPayload)

	case "textDocument/didSave":
		castedPayload, ok := payload.(*lsp.DidSaveTextDocumentParams)
		if !ok {
			return nil, ErrBadPayloadType
		}
		return nil, s.DidSave(ctx, castedPayload)

	case "textDocument/documentColor":
		castedPayload, ok := payload.(*lsp.DocumentColorParams)
		if !ok {
			return nil, ErrBadPayloadType
		}
		return s.DocumentColor(ctx, castedPayload)

	case "textDocument/documentHighlight":
		castedPayload, ok := payload.(*lsp.DocumentHighlightParams)
		if !ok {
			return nil, ErrBadPayloadType
		}
		return s.DocumentHighlight(ctx, castedPayload)

	case "textDocument/documentLink":
		castedPayload, ok := payload.(*lsp.DocumentLinkParams)
		if !ok {
			return nil, ErrBadPayloadType
		}
		return s.DocumentLink(ctx, castedPayload)

	case "textDocument/documentLinkResolve":
		castedPayload, ok := payload.(*lsp.DocumentLink)
		if !ok {
			return nil, ErrBadPayloadType
		}
		return s.DocumentLinkResolve(ctx, castedPayload)

	case "textDocument/documentSymbol":
		castedPayload, ok := payload.(*lsp.DocumentSymbolParams)
		if !ok {
			return nil, ErrBadPayloadType
		}
		return s.DocumentSymbol(ctx, castedPayload)

	case "textDocument/executeCommand":
		castedPayload, ok := payload.(*lsp.ExecuteCommandParams)
		if !ok {
			return nil, ErrBadPayloadType
		}
		return s.ExecuteCommand(ctx, castedPayload)

	case "exit":
		return nil, s.Exit(ctx)

	case "textDocument/foldingRanges":
		castedPayload, ok := payload.(*lsp.FoldingRangeParams)
		if !ok {
			return nil, ErrBadPayloadType
		}
		return s.FoldingRanges(ctx, castedPayload)

	case "textDocument/formatting":
		castedPayload, ok := payload.(*lsp.DocumentFormattingParams)
		if !ok {
			return nil, ErrBadPayloadType
		}
		return s.Formatting(ctx, castedPayload)

	case "textDocument/hover":
		castedPayload, ok := payload.(*lsp.HoverParams)
		if !ok {
			return nil, ErrBadPayloadType
		}
		return s.Hover(ctx, castedPayload)

	case "textDocument/implementation":
		castedPayload, ok := payload.(*lsp.ImplementationParams)
		if !ok {
			return nil, ErrBadPayloadType
		}
		return s.Implementation(ctx, castedPayload)

	case "textDocument/incomingCalls":
		castedPayload, ok := payload.(*lsp.CallHierarchyIncomingCallsParams)
		if !ok {
			return nil, ErrBadPayloadType
		}
		return s.IncomingCalls(ctx, castedPayload)

	case "initialize":
		castedPayload, ok := payload.(*lsp.InitializeParams)
		if !ok {
			return nil, ErrBadPayloadType
		}
		return s.Initialize(ctx, castedPayload)

	case "initialized":
		castedPayload, ok := payload.(*lsp.InitializedParams)
		if !ok {
			return nil, ErrBadPayloadType
		}
		return nil, s.Initialized(ctx, castedPayload)

	case "textDocument/linkedEditingRange":
		castedPayload, ok := payload.(*lsp.LinkedEditingRangeParams)
		if !ok {
			return nil, ErrBadPayloadType
		}
		return s.LinkedEditingRange(ctx, castedPayload)

	case "textDocument/logTrace":
		castedPayload, ok := payload.(*lsp.LogTraceParams)
		if !ok {
			return nil, ErrBadPayloadType
		}
		return nil, s.LogTrace(ctx, castedPayload)

	case "textDocument/moniker":
		castedPayload, ok := payload.(*lsp.MonikerParams)
		if !ok {
			return nil, ErrBadPayloadType
		}
		return s.Moniker(ctx, castedPayload)

	case "textDocument/onTypeFormatting":
		castedPayload, ok := payload.(*lsp.DocumentOnTypeFormattingParams)
		if !ok {
			return nil, ErrBadPayloadType
		}
		return s.OnTypeFormatting(ctx, castedPayload)

	case "textDocument/outgoingCalls":
		castedPayload, ok := payload.(*lsp.CallHierarchyOutgoingCallsParams)
		if !ok {
			return nil, ErrBadPayloadType
		}
		return s.OutgoingCalls(ctx, castedPayload)

	case "textDocument/prepareCallHierarchy":
		castedPayload, ok := payload.(*lsp.CallHierarchyPrepareParams)
		if !ok {
			return nil, ErrBadPayloadType
		}
		return s.PrepareCallHierarchy(ctx, castedPayload)

	case "textDocument/prepareRename":
		castedPayload, ok := payload.(*lsp.PrepareRenameParams)
		if !ok {
			return nil, ErrBadPayloadType
		}
		return s.PrepareRename(ctx, castedPayload)

	case "textDocument/rangeFormatting":
		castedPayload, ok := payload.(*lsp.DocumentRangeFormattingParams)
		if !ok {
			return nil, ErrBadPayloadType
		}
		return s.RangeFormatting(ctx, castedPayload)

	case "textDocument/references":
		castedPayload, ok := payload.(*lsp.ReferenceParams)
		if !ok {
			return nil, ErrBadPayloadType
		}
		return s.References(ctx, castedPayload)

	case "textDocument/rename":
		castedPayload, ok := payload.(*lsp.RenameParams)
		if !ok {
			return nil, ErrBadPayloadType
		}
		return s.Rename(ctx, castedPayload)

	case "textDocument/semanticTokens/full":
		castedPayload, ok := payload.(*lsp.SemanticTokensParams)
		if !ok {
			return nil, ErrBadPayloadType
		}
		return s.SemanticTokensFull(ctx, castedPayload)

	case "textDocument/semanticTokensFullDelta":
		castedPayload, ok := payload.(*lsp.SemanticTokensDeltaParams)
		if !ok {
			return nil, ErrBadPayloadType
		}
		return s.SemanticTokensFullDelta(ctx, castedPayload)

	case "textDocument/semanticTokensRange":
		castedPayload, ok := payload.(*lsp.SemanticTokensRangeParams)
		if !ok {
			return nil, ErrBadPayloadType
		}
		return s.SemanticTokensRange(ctx, castedPayload)

	case "textDocument/semanticTokensRefresh":
		return nil, s.SemanticTokensRefresh(ctx)

	case "textDocument/setTrace":
		castedPayload, ok := payload.(*lsp.SetTraceParams)
		if !ok {
			return nil, ErrBadPayloadType
		}
		return nil, s.SetTrace(ctx, castedPayload)

	case "textDocument/showDocument":
		castedPayload, ok := payload.(*lsp.ShowDocumentParams)
		if !ok {
			return nil, ErrBadPayloadType
		}
		return s.ShowDocument(ctx, castedPayload)

	case "shutdown":
		return nil, s.Shutdown(ctx)

	case "textDocument/signatureHelp":
		castedPayload, ok := payload.(*lsp.SignatureHelpParams)
		if !ok {
			return nil, ErrBadPayloadType
		}
		return s.SignatureHelp(ctx, castedPayload)

	case "textDocument/symbols":
		castedPayload, ok := payload.(*lsp.WorkspaceSymbolParams)
		if !ok {
			return nil, ErrBadPayloadType
		}
		return s.Symbols(ctx, castedPayload)

	case "textDocument/typeDefinition":
		castedPayload, ok := payload.(*lsp.TypeDefinitionParams)
		if !ok {
			return nil, ErrBadPayloadType
		}
		return s.TypeDefinition(ctx, castedPayload)

	case "textDocument/willCreateFiles":
		castedPayload, ok := payload.(*lsp.CreateFilesParams)
		if !ok {
			return nil, ErrBadPayloadType
		}
		return s.WillCreateFiles(ctx, castedPayload)

	case "textDocument/willDeleteFiles":
		castedPayload, ok := payload.(*lsp.DeleteFilesParams)
		if !ok {
			return nil, ErrBadPayloadType
		}
		return s.WillDeleteFiles(ctx, castedPayload)

	case "textDocument/willRenameFiles":
		castedPayload, ok := payload.(*lsp.RenameFilesParams)
		if !ok {
			return nil, ErrBadPayloadType
		}
		return s.WillRenameFiles(ctx, castedPayload)

	case "textDocument/willSave":
		castedPayload, ok := payload.(*lsp.WillSaveTextDocumentParams)
		if !ok {
			return nil, ErrBadPayloadType
		}
		return nil, s.WillSave(ctx, castedPayload)

	case "textDocument/willSaveWaitUntil":
		castedPayload, ok := payload.(*lsp.WillSaveTextDocumentParams)
		if !ok {
			return nil, ErrBadPayloadType
		}
		return s.WillSaveWaitUntil(ctx, castedPayload)

	case "textDocument/workDoneProgressCancel":
		castedPayload, ok := payload.(*lsp.WorkDoneProgressCancelParams)
		if !ok {
			return nil, ErrBadPayloadType
		}
		return nil, s.WorkDoneProgressCancel(ctx, castedPayload)

	}
	return nil, ErrUnknownMethod
}
