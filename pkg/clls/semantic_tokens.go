package clls

import (
	"fmt"
	"sort"
	"strings"
	"unicode"

	"github.com/clls-dev/clls/pkg/lsph"
	lsp "go.lsp.dev/protocol"
	"go.uber.org/zap"
)

func (m *Module) SemanticTokens(l *zap.Logger) ([]uint32, error) {
	inserts := []insert(nil)

	if m.IsMod {
		inserts = append(inserts, insert{Kind: "keyword", Token: m.ModToken})
		if m.Args != nil {
			inserts = insertParamsTokens(inserts, m.Args)
		}
	}

	if m.Constants != nil {
		for _, c := range m.Constants {
			if c.Token != nil {
				inserts = append(inserts, insert{Kind: "keyword", Token: c.Token})
			}
			if t, ok := c.Name.(*Token); ok && t != nil {
				inserts = append(inserts, insert{Kind: "variable", Modifiers: []lsp.SemanticTokenModifiers{lsp.SemanticTokenModifierReadonly}, Token: t})
			}
			inserts = insertBody(inserts, c.Value, BuiltinFuncsByName)
		}
	}

	for _, t := range m.Comments {
		inserts = append(inserts, insert{Kind: "comment", Token: t})
	}

	for _, incl := range m.Includes {
		inserts = append(inserts, insert{Kind: "keyword", Token: incl.Token})
		if t, ok := incl.Value.(*Token); ok && t != nil {
			inserts = append(inserts, insert{Kind: "string", Token: t})
		}
	}

	allFuncs := map[string]*Function{}
	for k, v := range BuiltinFuncsByName {
		allFuncs[k] = v
	}
	for k, v := range m.FunctionsByName {
		allFuncs[k] = v
	}
	if len(m.Functions) > 0 {
		for _, f := range m.Functions {
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

	if m.Main != nil {
		inserts = insertBody(inserts, m.Main, allFuncs)
	}

	data := []uint32(nil)

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
			tt, tm, err := tokenInfo(l, inserts[0].Kind, inserts[0].Modifiers, &StandardSemanticTokensLegend)
			if err != nil {
				panic(err)
			}
			ltoks = append(ltoks, lsph.SemanticToken{
				DeltaLine:      uint32(t.Line),
				DeltaStartChar: uint32(t.StartChar),
				Length:         uint32(len(t.Text)),
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

				tt, tm, err := tokenInfo(l, in.Kind, in.Modifiers, &StandardSemanticTokensLegend)
				if err != nil {
					panic(err)
				}
				ltoks = append(ltoks, lsph.SemanticToken{
					DeltaLine:      uint32(deltaLine),
					DeltaStartChar: uint32(deltaStartChar),
					Length:         uint32(len(t.Text)),
					TokenType:      tt,
					TokenModifiers: tm,
				})
			}

			data = ltoks.Flat()
		}
	}

	return data, nil
}

var StandardSemanticTokenTypes = []lsp.SemanticTokenTypes{
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

var StandardSemanticTokenModifiers = []lsp.SemanticTokenModifiers{
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

var StandardSemanticTokensLegend = lsp.SemanticTokensLegend{
	TokenTypes:     StandardSemanticTokenTypes,
	TokenModifiers: StandardSemanticTokenModifiers,
}

func tokenInfo(l *zap.Logger, kind lsp.SemanticTokenTypes, mods []lsp.SemanticTokenModifiers, legend *lsp.SemanticTokensLegend) (uint32, uint32, error) {
	tt := -1
	for i, v := range legend.TokenTypes {
		if v == kind {
			tt = i
			break
		}
	}
	if tt == -1 {
		return 0, 0, fmt.Errorf("unknown token type '%s'", kind)
	}
	tm := uint32(0)
	for _, m := range mods {
		mv := -1
		for i, v := range legend.TokenModifiers {
			if v == m {
				mv = i
				break
			}
		}
		if mv == -1 {
			return 0, 0, fmt.Errorf("unknown token modifier '%s'", m)
		}
		tm |= 1 << uint32(mv)
	}
	return uint32(tt), tm, nil
}

func insertParamsTokens(inserts []insert, a interface{}) []insert {
	switch a := a.(type) {
	case *Token:
		if a.Value != "." {
			inserts = append(inserts, insert{Kind: "parameter", Modifiers: []lsp.SemanticTokenModifiers{lsp.SemanticTokenModifierReadonly}, Token: a})
		}
	case *ASTNode:
		for _, ac := range a.Children {
			inserts = insertParamsTokens(inserts, ac)
		}
	}
	return inserts
}

func insertBody(inserts []insert, node *CodeBody, funcsByName map[string]*Function) []insert {
	if node == nil {
		return inserts
	}
	switch node.Kind {
	case IfBodyKind:
		if node.Token != nil {
			inserts = append(inserts, insert{Kind: "keyword", Token: node.Token})
		}
		inserts = insertBody(inserts, node.IfCond, funcsByName)
		inserts = insertBody(inserts, node.IfBranch, funcsByName)
		inserts = insertBody(inserts, node.ElseBranch, funcsByName)
	case CallBodyKind:
		kind := lsp.SemanticTokenFunction
		mods := []lsp.SemanticTokenModifiers(nil)
		if fn, ok := funcsByName[node.Function.Value]; ok && fn.Builtin {
			switch fn.Name.Value {
			case "x":
				kind = lsp.SemanticTokenKeyword
			default:
				mods = append(mods, lsp.SemanticTokenModifierDefaultLibrary)
			}
		}
		if node.Token != nil {
			inserts = append(inserts, insert{Kind: kind, Modifiers: mods, Token: node.Token})
		}
		for _, a := range node.CallArgs {
			inserts = insertBody(inserts, a, funcsByName)
		}
	case OperatorBodyKind:
		if node.Token != nil {
			inserts = append(inserts, insert{Kind: lsp.SemanticTokenOperator, Token: node.Token})
		}
		for _, child := range node.Children {
			inserts = insertBody(inserts, child, funcsByName)
		}
	case ConstBodyKind:
		if node.Token != nil {
			inserts = append(inserts, insert{Kind: lsp.SemanticTokenVariable, Modifiers: []lsp.SemanticTokenModifiers{lsp.SemanticTokenModifierReadonly}, Token: node.Token})
		}
	case VarBodyKind:
		if node.Token != nil {
			inserts = append(inserts, insert{Kind: lsp.SemanticTokenParameter, Token: node.Token})
		}
	case FuncVarBodyKind:
		k := lsp.SemanticTokenFunction
		mods := []lsp.SemanticTokenModifiers(nil)
		if fn, ok := funcsByName[node.Function.Value]; ok && fn.Builtin {
			mods = append(mods, lsp.SemanticTokenModifierDefaultLibrary)
		}
		if node.Token != nil {
			inserts = append(inserts, insert{Kind: k, Modifiers: mods, Token: node.Token})
		}
	default:
		if node.Token != nil && node.Token.Value != "." {
			kind := lsp.SemanticTokenString
			numberStr := node.Token.Value
			if len(numberStr) > 0 && numberStr[0] == '-' {
				numberStr = numberStr[1:]
			}
			if i := strings.IndexFunc(numberStr, func(r rune) bool { return !unicode.IsNumber(r) }); i == -1 {
				kind = lsp.SemanticTokenNumber
			}
			inserts = append(inserts, insert{Kind: kind, Token: node.Token})
		}
		for _, child := range node.Children {
			inserts = insertBody(inserts, child, funcsByName)
		}
	}
	return inserts
}

type insert struct {
	Kind      lsp.SemanticTokenTypes
	Modifiers []lsp.SemanticTokenModifiers
	Token     *Token
}
