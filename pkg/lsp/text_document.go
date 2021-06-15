package lsp

/**
 * Defines how the host (editor) should sync document changes to the language
 * server.
 */
type TextDocumentSyncKind int

/**
 * Documents should not be synced at all.
 */
const (
	None = TextDocumentSyncKind(0)

	/**
	 * Documents are synced by always sending the full content
	 * of the document.
	 */
	Full = TextDocumentSyncKind(1)

	/**
	 * Documents are synced by sending the full content on open.
	 * After that only incremental updates to the document are
	 * send.
	 */
	Incremental = TextDocumentSyncKind(2)
)

type TextDocumentSyncOptions struct {
	/**
	 * Open and close notifications are sent to the server. If omitted open
	 * close notification should not be sent.
	 */
	OpenClose bool `json:"openClose,omitempty"`

	/**
	 * Change notifications are sent to the server. See
	 * TextDocumentSyncKind.None, TextDocumentSyncKind.Full and
	 * TextDocumentSyncKind.Incremental. If omitted it defaults to
	 * TextDocumentSyncKind.None.
	 */
	Change TextDocumentSyncKind `json:"change,omitempty"`
}

type SemanticTokens struct {
	/*
		An optional result id. If provided and clients support delta updating
		the client will include the result id in the next semantic token request.
		A server can then instead of computing all semantic tokens again simply
		send a delta.
	*/
	ResultID string `json:"resultId,omitempty"`

	/*
		The actual tokens.
	*/
	Data []UInteger `json:"data"`
}

type WorkspaceEdit struct {
	/**
	 * Holds changes to existing resources.
	 */
	Changes map[DocumentUri][]TextEdit `json:"changes,omitempty"`

	/**
	 * Depending on the client capability
	 * `workspace.workspaceEdit.resourceOperations` document changes are either
	 * an array of `TextDocumentEdit`s to express changes to n different text
	 * documents where each text document edit addresses a specific version of
	 * a text document. Or it can contain above `TextDocumentEdit`s mixed with
	 * create, rename and delete file / folder operations.
	 *
	 * Whether a client supports versioned document edits is expressed via
	 * `workspace.workspaceEdit.documentChanges` client capability.
	 *
	 * If a client neither supports `documentChanges` nor
	 * `workspace.workspaceEdit.resourceOperations` then only plain `TextEdit`s
	 * using the `changes` property are supported.
	 */
	/*documentChanges?: (
		TextDocumentEdit[] |
		(TextDocumentEdit | CreateFile | RenameFile | DeleteFile)[]
	);*

	/**
	 * A map of change annotations that can be referenced in
	 * `AnnotatedTextEdit`s or create, rename and delete file / folder
	 * operations.
	 *
	 * Whether clients honor this property depends on the client capability
	 * `workspace.changeAnnotationSupport`.
	 *
	 * @since 3.16.0
	*/
	//changeAnnotations?: {
	//	[id: string /* ChangeAnnotationIdentifier */]: ChangeAnnotation;
	//};
}
