package lsp

type TextDocumentIdentifier struct {
	/*
	  The text document's URI.
	*/
	URI DocumentUri `json:"uri"`
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
	 * Trim trailing whitespace on a line.
	 *
	 * @since 3.15.0
	 */
	//trimTrailingWhitespace?: boolean;

	/**
	 * Insert a newline character at the end of the file if one does not exist.
	 *
	 * @since 3.15.0
	 */
	//insertFinalNewline?: boolean;

	/**
	 * Trim all newlines after the final newline at the end of the file.
	 *
	 * @since 3.15.0
	 */
	//trimFinalNewlines?: boolean;

	/**
	 * Signature for further properties.
	 */
	//[key: string]: boolean | integer | string;
}
