package lsp

type TextDocumentIdentifier struct {
	/*
	  The text document's URI.
	*/
	URI DocumentURI `json:"uri"`
}

type DocumentFormattingParams struct {
	/*
	  The document to format.
	*/
	TextDocument TextDocumentIdentifier `json:"textDocument"`

	/*
	  The format options.
	*/
	Options FormattingOptions `json:"options"`
}

type FormattingOptions struct {
	/*
	  Size of a tab in spaces.
	*/
	TabSize UInteger `json:"tabSize"`

	/*
	   Prefer spaces over tabs.
	*/
	InsertSpaces bool `json:"insertSpaces"`

	/**
	Trim trailing whitespace on a line.
	 *
	@since 3.15.0
	*/
	//trimTrailingWhitespace?: boolean;

	/**
	Insert a newline character at the end of the file if one does not exist.
	 *
	@since 3.15.0
	*/
	//insertFinalNewline?: boolean;

	/**
	Trim all newlines after the final newline at the end of the file.
	 *
	@since 3.15.0
	*/
	//trimFinalNewlines?: boolean;

	/**
	Signature for further properties.
	*/
	//[key: string]: boolean | integer | string;
}

type DidOpenTextDocumentParams struct {
	// The document that was opened.
	TextDocument TextDocumentItem `json:"textDocument"`
}

type DidCloseTextDocumentParams struct {
	// The document that was closed.
	TextDocument TextDocumentItem `json:"textDocument"`
}

type DidChangeTextDocumentParams struct {
	/*
		The document that did change. The version number points
		to the version after all provided content changes have
		been applied.
	*/
	TextDocument VersionedTextDocumentIdentifier `json:"textDocument"`

	/*
		The actual content changes. The content changes describe single state
		changes to the document. So if there are two content changes c1 (at
		array index 0) and c2 (at array index 1) for a document in state S then
		c1 moves the document from S to S' and c2 from S' to S''. So c1 is
		computed on the state S and c2 is computed on the state S'.

		To mirror the content of a document using change events use the following
		approach:

		- start with the same initial content

		- apply the 'textDocument/didChange' notifications in the order you receive them.

		- apply the `TextDocumentContentChangeEvent`s in a single notification in the order you receive them.
	*/
	ContentChanges []TextDocumentContentChangeEvent `json:"contentChanges"`
}

type VersionedTextDocumentIdentifier struct {
	TextDocumentIdentifier

	/**
	 * The version number of this document.
	 *
	 * The version number of a document will increase after each change,
	 * including undo/redo. The number doesn't need to be consecutive.
	 */
	Version Integer `json:"version"`
}

type TextDocumentContentChangeEvent struct {
	/**
	 * The range of the document that changed.
	 */
	Range *Range `json:"range,omitempty"`

	/**
	 * The optional length of the range that got replaced.
	 *
	 * @deprecated use range instead.
	 */
	RangeLength UInteger `json:"rangeLength,omitempty"`

	/**
	 * The new text for the provided range.
	 */
	Text string `json:"text"`
}

type TextDocumentItem struct {
	// The text document's URI.
	URI DocumentURI `json:"uri"`
	// The text document's language identifier.
	LanguageID string `json:"languageId"`
	// The version number of this document (it will increase after each change, including undo/redo).
	Version Integer `json:"version"`
	// The content of the opened text document.
	Text string `json:"text"`
}
