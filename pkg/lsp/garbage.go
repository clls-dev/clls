package lsp

type InitializedParams struct{}
type InitializeParams struct{}
type ShutdownParams struct{}
type ShutdownResult struct{}

type RenameParams struct {
	TextDocumentPositionParams
	WorkDoneProgressParams

	/*
		The new name of the symbol. If the given name is not valid the
		request must return a [ResponseError](#ResponseError) with an
		appropriate message set.
	*/
	NewName string `json:"newName"`
}

type TextDocumentPositionParams struct {
	/*
		The text document.
	*/
	TextDocument TextDocumentIdentifier `json:"textDocument"`

	/*
		The position inside the text document.
	*/
	Position Position `json:"position"`
}

type WorkDoneProgressParams struct {
	/*
		An optional token that a server can use to report work done progress.
	*/
	WorkDoneToken *ProgressToken `json:"workDoneToken,omitempty"`
}

type SemanticTokensParams struct {
	WorkDoneProgressParams
	PartialResultParams
	/*
		The text document.
	*/
	TextDocument TextDocumentIdentifier `json:"textDocument"`
}

type PartialResultParams struct {
	/*
		An optional token that a server can use to report partial results (e.g.
		streaming) to the client.
	*/
	PartialResultToken *ProgressToken `json:"textDocument,omitempty"`
}
