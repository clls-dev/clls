package lsp

// URI

type DocumentURI = string
type URI = string

// Regular Expression

type RegularExpressionsClientCapabilities struct {
	Engine  string  `json:"engine"`
	Version *string `json:"version,omitempty"`
}

// Text Documents

var EOL = []string{"\n", "\r\n", "\r"}

// Position

type Position struct {
	Line      UInteger `json:"line"`
	Character UInteger `json:"character"`
}

// Range

type Range struct {
	Start Position `json:"start"`
	End   Position `json:"end"`
}

// Location

type Location struct {
	URI   DocumentURI `json:"uri"`
	Range Range       `json:"range"`
}

// LocationLink

type LocationLink struct {
	OriginSelectionRange *Range      `json:"originSelectionRange"`
	TargetURI            DocumentURI `json:"targetUri"`
	TargetRange          Range       `json:"targetRange"`
	TargetSelectionRange Range       `json:"targetSelectionRange"`
}

// Diagnostic

type Diagnostic struct {
	/**
	 * The range at which the message applies.
	 */
	Range Range

	/**
	 * The diagnostic's severity. Can be omitted. If omitted it is up to the
	 * client to interpret diagnostics as error, warning, info or hint.
	 */
	Severity *DiagnosticSeverity

	/**
	 * The diagnostic's code, which might appear in the user interface.
	 */
	Code *IntegerOrString

	/**
	 * An optional property to describe the error code.
	 *
	 * @since 3.16.0
	 */
	CodeDescription *CodeDescription

	/**
	 * A human-readable string describing the source of this
	 * diagnostic, e.g. 'typescript' or 'super lint'.
	 */
	Source *string

	/**
	 * The diagnostic's message.
	 */
	Message string

	/**
	 * Additional metadata about the diagnostic.
	 *
	 * @since 3.15.0
	 */
	Tags *[]DiagnosticTag

	/**
	 * An array of related diagnostic information, e.g. when symbol-names within
	 * a scope collide all definitions can be marked via this property.
	 */
	RelatedInformation *[]DiagnosticRelatedInformation

	/*
		A data entry field that is preserved between a
		`textDocument/publishDiagnostics` notification and
		`textDocument/codeAction` request.

		@since 3.16.0
	*/
	Data interface{}
}

type DiagnosticSeverity int

const (
	/*
		Reports an error.
	*/
	Error = DiagnosticSeverity(1)
	/*
		Reports a warning.
	*/
	Warning = DiagnosticSeverity(2)
	/*
		Reports an information.
	*/
	Information = DiagnosticSeverity(3)
	/*
		Reports a hint.
	*/
	Hint = DiagnosticSeverity(4)
)

/*
The diagnostic tags.

@since 3.15.0
*/
type DiagnosticTag int

const (
	/**
	* Unused or unnecessary code.
	*
	* Clients are allowed to render diagnostics with this tag faded out
	* instead of having an error squiggle.
	 */
	Unnecessary = DiagnosticTag(1)
	/**
	* Deprecated or obsolete code.
	*
	* Clients are allowed to rendered diagnostics with this tag strike through.
	 */
	Deprecated = DiagnosticTag(2)
)

type DiagnosticRelatedInformation struct {
	Location Location
	Message  string
}

type CodeDescription struct {
	HRef URI
}

// Command

type Command struct {
	Title     string
	Command   string
	Arguments []interface{}
}

// TextEdit & AnnotatedTextEdit

type TextEdit struct {
	Range   Range
	NewText string
}

type ChangeAnnotation struct {
	Label             string
	NeedsConfirmation bool `json:"needsConfirmation,omitempty"`
	Description       string
}

type ChangeAnnotationIdentifier = string

type AnnotatedTextEdit struct {
	TextEdit
	AnnotationID ChangeAnnotationIdentifier
}

// TextEdit

// Text Document Edit

type TextDocumentEdit struct {
	//TextDocument OptionalVersionedTextDocumentIdentifier
	Edits []interface{} //  (TextEdit | AnnotatedTextEdit)[];
}

// File Resource changes

type CreateFileOptions struct {
	Overwrite      bool `json:"overwrite,omitempty"`
	IgnoreIfExists bool `json:"ignoreIfExists,omitempty"`
}
