package clls

import (
	"fmt"
	"sort"
	"strings"
	"unicode"

	"github.com/clls-dev/clls/pkg/lsp"
	"github.com/clls-dev/clls/pkg/lsph"
	"go.uber.org/zap"
)

func (m *Module) SemanticTokens(l *zap.Logger) ([]lsp.UInteger, error) {
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
				inserts = append(inserts, insert{Kind: "variable", Modifiers: []string{"readonly"}, Token: t})
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

	data := []lsp.UInteger(nil)

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
			tt, tm, err := tokenInfo(l, inserts[0].Kind, inserts[0].Modifiers, &lsp.StandardSemanticTokensLegend)
			if err != nil {
				panic(err)
			}
			ltoks = append(ltoks, lsph.SemanticToken{
				DeltaLine:      lsp.UInteger(t.Line),
				DeltaStartChar: lsp.UInteger(t.StartChar),
				Length:         lsp.UInteger(len(t.Text)),
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

				tt, tm, err := tokenInfo(l, in.Kind, in.Modifiers, &lsp.StandardSemanticTokensLegend)
				if err != nil {
					panic(err)
				}
				ltoks = append(ltoks, lsph.SemanticToken{
					DeltaLine:      lsp.UInteger(deltaLine),
					DeltaStartChar: lsp.UInteger(deltaStartChar),
					Length:         lsp.UInteger(len(t.Text)),
					TokenType:      tt,
					TokenModifiers: tm,
				})
			}

			data = ltoks.Flat()
		}
	}

	return data, nil
}

func tokenInfo(l *zap.Logger, kind string, mods []string, legend *lsp.SemanticTokensLegend) (lsp.UInteger, lsp.UInteger, error) {
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
	tm := lsp.UInteger(0)
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
		tm |= 1 << lsp.UInteger(mv)
	}
	return lsp.UInteger(tt), tm, nil
}

func insertParamsTokens(inserts []insert, a interface{}) []insert {
	switch a := a.(type) {
	case *Token:
		if a.Value != "." {
			inserts = append(inserts, insert{Kind: "parameter", Modifiers: []string{"readonly"}, Token: a})
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
		kind := "function"
		mods := []string(nil)
		if fn, ok := funcsByName[node.Function.Value]; ok && fn.Builtin {
			switch fn.Name.Value {
			case "x":
				kind = "keyword"
			default:
				mods = append(mods, "defaultLibrary")
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
			inserts = append(inserts, insert{Kind: "operator", Token: node.Token})
		}
		for _, child := range node.Children {
			inserts = insertBody(inserts, child, funcsByName)
		}
	case ConstBodyKind:
		if node.Token != nil {
			inserts = append(inserts, insert{Kind: "variable", Modifiers: []string{"readonly"}, Token: node.Token})
		}
	case VarBodyKind:
		if node.Token != nil {
			inserts = append(inserts, insert{Kind: "parameter", Token: node.Token})
		}
	case FuncVarBodyKind:
		k := "function"
		mods := []string(nil)
		if fn, ok := funcsByName[node.Function.Value]; ok && fn.Builtin {
			k = "function"
			mods = append(mods, "defaultLibrary")
		}
		if node.Token != nil {
			inserts = append(inserts, insert{Kind: k, Modifiers: mods, Token: node.Token})
		}
	default:
		if node.Token != nil && node.Token.Value != "." {
			kind := "string"
			numberStr := node.Token.Value
			if len(numberStr) > 0 && numberStr[0] == '-' {
				numberStr = numberStr[1:]
			}
			if i := strings.IndexFunc(numberStr, func(r rune) bool { return !unicode.IsNumber(r) }); i == -1 {
				kind = "number"
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
	Kind      string
	Modifiers []string
	Token     *Token
}
