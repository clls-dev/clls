package clls

import (
	"strings"

	"github.com/pkg/errors"
	lsp "go.lsp.dev/protocol"
	"go.uber.org/zap"
)

func Prettify(l *zap.Logger, fo *lsp.FormattingOptions, text string, documentURI lsp.DocumentURI) (string, int, error) {
	tokens, err := tokenizeSync(text, documentURI)
	if err != nil {
		return "", 0, errors.Wrap(err, "tokenize")
	}
	//l.Debug("got tokens", zap.Any("tokens", tokens))
	final := ""
	indent := "\t"
	if fo.InsertSpaces {
		indent = strings.Repeat(" ", int(fo.TabSize))
	}
	level := 0
	linesCount := 1
	l.Debug("start prettify")
	for i := 0; i < len(tokens); i++ {
		t := tokens[i]
		pt := (*Token)(nil)
		if i > 0 {
			pt = tokens[i-1]
		}
		if t.Kind == lineReturnToken {
			linesCount++
			final += t.Text
			continue
		}

		insertIndent := false
		if pt != nil && pt.Kind == lineReturnToken {
			for ; i < len(tokens); i++ {
				t = tokens[i]
				if t.Kind != spaceToken {
					break
				}
			}
			if i == len(tokens) {
				break
			}
			insertIndent = true
		}
		if t.Kind == parensCloseToken {
			if level > 0 {
				level--
			}
		}
		if insertIndent {
			final += strings.Repeat(indent, level)
		}
		final += t.Text
		if t.Kind == parensOpenToken {
			level++
		}
	}

	return final, linesCount, nil
}
