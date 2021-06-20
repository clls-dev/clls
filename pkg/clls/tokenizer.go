package clls

import (
	"encoding/json"
	"fmt"
	"regexp"
	"unicode"

	"github.com/pkg/errors"
	lsp "go.lsp.dev/protocol"
)

type Token struct {
	Value       string
	Index       int
	Kind        tokenKind
	Text        string
	Line        int
	StartChar   int // relative to line start
	DocumentURI lsp.DocumentURI
}

func (t *Token) EndChar() int {
	if t.Kind == lineReturnToken {
		return 0
	}
	return t.StartChar + len(t.Text)
}

func (t *Token) EndLine() int {
	if t.Kind == lineReturnToken {
		return t.Line + 1
	}
	return t.Line
}

func (t *Token) Range() lsp.Range {
	return lsp.Range{
		Start: lsp.Position{Line: uint32(t.Line), Character: uint32(t.StartChar)},
		End:   lsp.Position{Line: uint32(t.EndLine()), Character: uint32(t.EndChar())},
	}
}

type tokenKind int

const (
	unknownToken = tokenKind(iota)
	basicToken
	parensOpenToken
	parensCloseToken
	quoteToken
	commentToken
	spaceToken
	lineReturnToken
)

func tokenizeSync(text string, documentURI lsp.DocumentURI) ([]*Token, error) {
	tch, errptr := tokenize(text, documentURI)
	tokens := []*Token(nil)
	for token := range tch {
		tokens = append(tokens, token)
	}
	if *errptr != nil {
		return nil, *errptr
	}
	return tokens, nil
}

func tokenize(text string, documentURI lsp.DocumentURI) (_ chan *Token, errptr *error) {
	err := error(nil)
	errptr = &err
	ch := make(chan *Token)
	go func() {
		defer close(ch)
		wordStart := -1
		i := 0
		lines := []int{0}
		updateTokenLine := func(t *Token) *Token {
			found := false
			for i := 1; i < len(lines); i++ {
				l := lines[i]
				if t.Index < l {
					t.Line = i - 1
					t.StartChar = t.Index - lines[i-1]
					found = true
					break
				}
			}
			if !found {
				t.Line = len(lines) - 1
				t.StartChar = t.Index - lines[len(lines)-1]
			}
			return t
		}
		cutWord := func() {
			if wordStart != -1 && i-wordStart > 0 {
				ch <- updateTokenLine(&Token{
					Value:       text[wordStart:i],
					Index:       wordStart,
					Text:        text[wordStart:i],
					Kind:        basicToken,
					DocumentURI: documentURI,
				})
				wordStart = -1
			}
		}
		for i < len(text) {
			c := text[i]
			if c == ';' {
				cutWord()
				commentStart := i
				commentEnd := regexp.MustCompile(`\r?\n`).FindIndex([]byte(text[i:]))[0] + i
				if commentEnd == -1 {
					commentEnd = len(text)
				}
				ch <- updateTokenLine(&Token{
					Value:       text[commentStart+1 : commentEnd],
					Index:       commentStart,
					Text:        text[commentStart:commentEnd],
					Kind:        commentToken,
					DocumentURI: documentURI,
				})
				i = commentEnd
			} else if c == '(' || c == ')' {
				cutWord()
				kind := parensOpenToken
				if c == ')' {
					kind = parensCloseToken
				}
				ch <- updateTokenLine(&Token{
					Value:       string(c),
					Index:       i,
					Text:        string(c),
					Kind:        kind,
					DocumentURI: documentURI,
				})
				i++
			} else if is := regexp.MustCompile(`^\r?\n`).Find([]byte(text[i:])); len(is) > 0 {
				cutWord()
				lines = append(lines, i+len(is))
				ch <- updateTokenLine(&Token{
					Value:       text[i : i+len(is)],
					Index:       i,
					Text:        text[i : i+len(is)],
					Kind:        lineReturnToken,
					DocumentURI: documentURI,
				})
				i += len(is)
			} else if unicode.IsSpace(rune(c)) {
				cutWord()
				spaceStart := i
				nextNonSpace := i
				for nextNonSpace < len(text) {
					if !unicode.IsSpace(rune(text[nextNonSpace])) {
						break
					}
					if is := regexp.MustCompile(`^\r?\n`).Find([]byte(text[nextNonSpace:])); len(is) > 0 {
						break
					}
					nextNonSpace++
				}
				ch <- updateTokenLine(&Token{
					Value:       text[spaceStart:nextNonSpace],
					Index:       spaceStart,
					Text:        text[spaceStart:nextNonSpace],
					Kind:        spaceToken,
					DocumentURI: documentURI,
				})
				i = nextNonSpace
			} else if c == '"' {
				cutWord()
				quoteStart := i
				quoteValue := ""
				j := i + 1
				for j < len(text) {
					c := text[j]

					if j < len(text)-1 && c == '\\' && text[j+1] == '"' {
						quoteValue += `"`
						j += 2
						continue
					}

					if c == '"' {
						ch <- updateTokenLine(&Token{
							Value:       quoteValue,
							Index:       quoteStart,
							Text:        text[quoteStart : j+1],
							Kind:        quoteToken,
							DocumentURI: documentURI,
						})
						quoteStart = -1
						quoteValue = ""
						break
					}

					quoteValue += string(c)
					j++
				}
				if j == len(text) {
					*errptr = errors.New("unclosed quote")
					return
				}
				i = j + 1
			} else {
				if wordStart == -1 {
					wordStart = i
				}
				i++
			}
		}
		cutWord()
	}()

	return ch, errptr
}

var tokenKindNames = map[tokenKind]string{
	basicToken:       "basic",
	parensCloseToken: "close",
	parensOpenToken:  "open",
	quoteToken:       "quote",
	commentToken:     "comment",
	spaceToken:       "space",
	lineReturnToken:  "line-return",
}

var tokenKindsByName = func() map[string]tokenKind {
	m := map[string]tokenKind{}
	for k, v := range tokenKindNames {
		m[v] = k
	}
	return m
}()

func (tk *tokenKind) UnmarshalJSON(b []byte) error {
	var v string
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	s, ok := tokenKindsByName[v]
	if !ok {
		return fmt.Errorf("unknown token kind '%s'", v)
	}
	*tk = s
	return nil
}

func (tk *tokenKind) MarshalJSON() ([]byte, error) {
	return json.Marshal(tk.String())
}
