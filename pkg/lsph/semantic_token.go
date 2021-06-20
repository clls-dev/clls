package lsph

type SemanticToken struct {
	DeltaLine      uint32
	DeltaStartChar uint32
	Length         uint32
	TokenType      uint32
	TokenModifiers uint32
}

func (l *SemanticToken) Slice() []uint32 {
	return []uint32{l.DeltaLine, l.DeltaStartChar, l.Length, l.TokenType, l.TokenModifiers}
}

type SemanticTokenSlice []SemanticToken

func (ls SemanticTokenSlice) Flat() []uint32 {
	lu := make([]uint32, len(ls)*5)
	for i, l := range ls {
		copy(lu[i*5:], l.Slice())
	}
	return lu
}
