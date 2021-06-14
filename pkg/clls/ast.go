package clls

import (
	"github.com/pkg/errors"
)

type ASTNode struct {
	Children   []interface{}
	OpenToken  *Token
	CloseToken *Token
}

func parseAST(tokch chan *Token) (*ASTNode, error) {
	parents := []*ASTNode{{}} // start with empty root
	for {
		t, ok := <-tokch
		if !ok {
			break
		}
		current := parents[len(parents)-1]
		switch t.Kind {
		case quoteToken, basicToken:
			current.Children = append(current.Children, t)
		case parensOpenToken:
			child := &ASTNode{}
			parents = append(parents, child)
			child.OpenToken = t
			current.Children = append(current.Children, child)
		case parensCloseToken:
			if len(parents) == 1 {
				continue
			}
			parents = parents[:len(parents)-1]
			parents[len(parents)-1].CloseToken = t
		}
	}
	if len(parents) < 1 {
		return nil, errors.New("unexpected internal state")
	}
	return parents[0], nil
}
