package lsp

type ServerCapabilities struct {
	/*
		The server provides semantic tokens support.

		Type must be SemanticTokensOptions or SemanticTokensRegistrationOptions
	*/
	SemanticTokensProvider interface{} `json:"semanticTokenProvider"`

	/*
		Defines how text documents are synced. Is either a detailed structure
		defining each notification or for backwards compatibility the
		TextDocumentSyncKind number. If omitted it defaults to
		`TextDocumentSyncKind.None`.

		Type must be TextDocumentSyncOptions or TextDocumentSyncKind
	*/
	TextDocumentSync interface{} `json:"textDocumentSync"`

	DeclarationProvider interface{} `json:"declarationProvider,omitempty"`

	/*
		The server provides completion support.
	*/
	CompletionProvider *CompletionOptions `json:"completionOptions,omitempty"`

	/*
		The server provides document formatting.

		Type must be nil, bool or DocumentFormattingOptions
	*/
	DocumentFormattingProvider interface{} `json:"documentFormattingProvider,omitempty"`

	/*
	  The server provides rename support. RenameOptions may only be
	  specified if the client states that it supports
	  `prepareSupport` in its initial `initialize` request.

	  Type must be nil, bool or RenameOptions
	*/
	RenameProvider interface{} `json:"renameProvider,omitempty"`

	/*
	  The server provides document symbol support.

	  Type must be nil, bool or DocumentSymbolOptions
	*/
	DocumentSymbolProvider interface{} `json:"documentSymbolProvider,omitempty"`

	/*
		Blabla
	*/
	ReferencesProvider interface{} `json:"referencesProvider,omitempty"`

	/*
	 Blabla
	*/
	DocumentHighlightProvider interface{} `json:"documentHighlightProvider,omitempty"`
}

type DocumentFormattingOptions struct {
}

type CompletionOptions struct {
	/**
	 * Most tools trigger completion request automatically without explicitly
	 * requesting it using a keyboard shortcut (e.g. Ctrl+Space). Typically they
	 * do so when the user starts to type an identifier. For example if the user
	 * types `c` in a JavaScript file code complete will automatically pop up
	 * present `console` besides others as a completion item. Characters that
	 * make up identifiers don't need to be listed here.
	 *
	 * If code complete should automatically be trigger on characters not being
	 * valid inside an identifier (for example `.` in JavaScript) list them in
	 * `triggerCharacters`.
	 */
	TriggerCharacters []string `json:"triggerCharacters"`

	/**
	 * The list of all possible characters that commit a completion. This field
	 * can be used if clients don't support individual commit characters per
	 * completion item. See client capability
	 * `completion.completionItem.commitCharactersSupport`.
	 *
	 * If a server provides both `allCommitCharacters` and commit characters on
	 * an individual completion item the ones on the completion item win.
	 *
	 * @since 3.2.0
	 */
	AllCommitCharacters []string `json:"allCommitCharacters"`

	/**
	 * The server provides support to resolve additional
	 * information for a completion item.
	 */
	ResolveProvider bool `json:"resolveProvider"`

	/**
	 * The server supports the following `CompletionItem` specific
	 * capabilities.
	 *
	 * @since 3.17.0 - proposed state
	 */
	CompletionItem *CompletionItem `json:"completionItem"`
}

type CompletionItem struct {
	/**
	 * The server has support for completion item label
	 * details (see also `CompletionItemLabelDetails`) when receiving
	 * a completion item in a resolve call.
	 *
	 * @since 3.17.0 - proposed state
	 */
	LabelDetailsSupport bool `json:"labelDetailsSupport"`
}

type SemanticTokensOptions struct {
	// The legend used by the server
	Legend SemanticTokensLegend `json:"legend"`
	// Server supports providing semantic tokens for a specific range of a document.
	Range bool `json:"range,omitempty"`
	// Server supports providing semantic tokens for a full document.
	Full interface{} `json:"full,omitempty"`
}

type SemanticTokensRegistrationOptions struct {
	SemanticTokensOptions
	StaticRegistrationOptions
	TextDocumentRegistrationOptions
}

type StaticRegistrationOptions struct {
	ID string `json:"id,omitempty"`
}

type TextDocumentRegistrationOptions struct {
	DocumentSelector DocumentSelector `json:"documentSelector,omitempty"`
}

type DocumentSelector = []DocumentFilter

type DocumentFilter struct {
	Language string `json:"language,omitempty"`
	Scheme   string `json:"scheme,omitempty"`
	Pattern  string `json:"pattern,omitempty"`
}

type SemanticTokensLegend struct {
	// The token types a server uses.
	TokenTypes []string `json:"tokenTypes"`
	// The token modifiers a server uses.
	TokenModifiers []string `json:"tokenModifiers"`
}

var SemanticTokenTypes = []string{
	"namespace",
	/**
	 * Represents a generic type. Acts as a fallback for types which
	 * can't be mapped to a specific type like class or enum.
	 */
	"type",
	"class",
	"enum",
	"interface",
	"struct",
	"typeParameter",
	"parameter",
	"variable",
	"property",
	"enumMember",
	"event",
	"function",
	"method",
	"macro",
	"keyword",
	"modifier",
	"comment",
	"string",
	"number",
	"regexp",
	"operator",
}

var SemanticTokenModifiers = []string{
	"declaration",
	"definition",
	"readonly",
	"static",
	"deprecated",
	"abstract",
	"async",
	"modification",
	"documentation",
	"defaultLibrary",
}

var StandardSemanticTokensLegend = SemanticTokensLegend{
	TokenTypes:     SemanticTokenTypes,
	TokenModifiers: SemanticTokenModifiers,
}
