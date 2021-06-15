package lsp

// Reply to document symbol must be DocumentSymbol[] | SymbolInformation[] | null

/*
  Represents information about programming constructs like variables, classes,
  interfaces etc.
*/
type SymbolInformation struct {
	/**
	 * The name of this symbol.
	 */
	//name: string;

	/**
	 * The kind of this symbol.
	 */
	//kind: SymbolKind;

	/**
	 * Tags for this symbol.
	 *
	 * @since 3.16.0
	 */
	//tags?: SymbolTag[];

	/**
	 * Indicates if this symbol is deprecated.
	 *
	 * @deprecated Use tags instead
	 */
	//deprecated?: boolean;

	/**
	 * The location of this symbol. The location's range is used by a tool
	 * to reveal the location in the editor. If the symbol is selected in the
	 * tool the range's start information is used to position the cursor. So
	 * the range usually spans more then the actual symbol's name and does
	 * normally include things like visibility modifiers.
	 *
	 * The range doesn't have to denote a node range in the sense of a abstract
	 * syntax tree. It can therefore not be used to re-construct a hierarchy of
	 * the symbols.
	 */
	//location: Location;

	/**
	 * The name of the symbol containing this symbol. This information is for
	 * user interface purposes (e.g. to render a qualifier in the user interface
	 * if necessary). It can't be used to re-infer a hierarchy for the document
	 * symbols.
	 */
	//containerName?: string;
}

/*
	Represents programming constructs like variables, classes, interfaces etc.
	that appear in a document. Document symbols can be hierarchical and they
	have two ranges: one that encloses its definition and one that points to its
	most interesting range, e.g. the range of an identifier.
*/
type DocumentSymbol struct {
	/*
		The name of this symbol. Will be displayed in the user interface and
		therefore must not be an empty string or a string only consisting of
		white spaces.
	*/
	Name string `json:"name"`

	/*
		More detail for this symbol, e.g the signature of a function.
	*/
	Detail string `json:"detail,omitempty"`

	/*
		The kind of this symbol.
	*/
	Kind SymbolKind `json:"kind"`

	/*
		Tags for this document symbol.

		@since 3.16.0
	*/
	Tags []SymbolTag `json:"tags,omitempty"`

	/*
		Indicates if this symbol is deprecated.

		@deprecated Use tags instead
	*/
	Deprecated bool `json:"deprecated,omitempty"`

	/*
		The range enclosing this symbol not including leading/trailing whitespace
		but everything else like comments. This information is typically used to
		determine if the clients cursor is inside the symbol to reveal in the
		symbol in the UI.
	*/
	Range Range `json:"range"`

	/*
		The range that should be selected and revealed when this symbol is being
		picked, e.g. the name of a function. Must be contained by the `range`.
	*/
	SelectionRange Range `json:"selectionRange"`

	/*
		Children of this symbol, e.g. properties of a class.
	*/
	Children []DocumentSymbol `json:"children,omitempty"`
}

/*
  A symbol kind.
*/
type SymbolKind int

const (
	SymbolKindFile          = SymbolKind(1)
	SymbolKindModule        = SymbolKind(2)
	SymbolKindNamespace     = SymbolKind(3)
	SymbolKindPackage       = SymbolKind(4)
	SymbolKindClass         = SymbolKind(5)
	SymbolKindMethod        = SymbolKind(6)
	SymbolKindProperty      = SymbolKind(7)
	SymbolKindField         = SymbolKind(8)
	SymbolKindConstructor   = SymbolKind(9)
	SymbolKindEnum          = SymbolKind(10)
	SymbolKindInterface     = SymbolKind(11)
	SymbolKindFunction      = SymbolKind(12)
	SymbolKindVariable      = SymbolKind(13)
	SymbolKindConstant      = SymbolKind(14)
	SymbolKindString        = SymbolKind(15)
	SymbolKindNumber        = SymbolKind(16)
	SymbolKindBoolean       = SymbolKind(17)
	SymbolKindArray         = SymbolKind(18)
	SymbolKindObject        = SymbolKind(19)
	SymbolKindKey           = SymbolKind(20)
	SymbolKindNull          = SymbolKind(21)
	SymbolKindEnumMember    = SymbolKind(22)
	SymbolKindStruct        = SymbolKind(23)
	SymbolKindEvent         = SymbolKind(24)
	SymbolKindOperator      = SymbolKind(25)
	SymbolKindTypeParameter = SymbolKind(26)
)

/*
  Symbol tags are extra annotations that tweak the rendering of a symbol.

  @since 3.16
*/
type SymbolTag int

const (
	/*
	   Render a symbol as obsolete, usually using a strike-out.
	*/
	SymbolTagDeprecated = SymbolTag(1)
)
