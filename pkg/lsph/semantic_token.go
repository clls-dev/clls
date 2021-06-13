package lsph

import "github.com/clls-dev/clls/pkg/lsp"

type SemanticToken struct {
	DeltaLine      lsp.UInteger
	DeltaStartChar lsp.UInteger
	Length         lsp.UInteger
	TokenType      lsp.UInteger
	TokenModifiers lsp.UInteger
}

func (l *SemanticToken) Slice() []lsp.UInteger {
	return []lsp.UInteger{l.DeltaLine, l.DeltaStartChar, l.Length, l.TokenType, l.TokenModifiers}
}

type SemanticTokenSlice []SemanticToken

func (ls SemanticTokenSlice) Flat() []lsp.UInteger {
	lu := make([]lsp.UInteger, len(ls)*5)
	for i, l := range ls {
		copy(lu[i*5:], l.Slice())
	}
	return lu
}
