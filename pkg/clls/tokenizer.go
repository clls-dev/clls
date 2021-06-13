package clls

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/pkg/errors"
)

type Token struct {
	Value     string
	Index     int
	Kind      tokenKind
	Text      string
	Line      int
	StartChar int // relative to line start
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
)

func tokenize(text string) (chan *Token, chan error) {
	ch := make(chan *Token)
	errCh := make(chan error)

	go func() {
		defer close(ch)
		defer close(errCh)
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
					Value: text[wordStart:i],
					Index: wordStart,
					Text:  text[wordStart:i],
					Kind:  basicToken,
				})
				wordStart = -1
			}
		}
		for i < len(text) {
			c := text[i]
			if c == ';' { // ignore comment
				cutWord()
				commentStart := i
				commentEnd := strings.IndexRune(text[i+1:], '\n')
				if commentEnd == -1 {
					commentEnd = len(text)
				} else {
					commentEnd += i + 1
				}
				ch <- updateTokenLine(&Token{
					Value: text[commentStart+1 : commentEnd],
					Index: commentStart,
					Text:  text[commentStart:commentEnd],
					Kind:  commentToken,
				})
				i = commentEnd
			} else if c == '(' || c == ')' {
				cutWord()
				kind := parensOpenToken
				if c == ')' {
					kind = parensCloseToken
				}
				ch <- updateTokenLine(&Token{
					Value: string(c),
					Index: i,
					Text:  string(c),
					Kind:  kind,
				})
				i++
			} else if unicode.IsSpace(rune(c)) {
				cutWord()
				spaceStart := i

				nextNonSpace := i
				for nextNonSpace < len(text) {
					if !unicode.IsSpace(rune(text[nextNonSpace])) {
						break
					}
					if is := regexp.MustCompile(`^\r?\n`).Find([]byte(text[nextNonSpace:])); len(is) > 0 {
						nextNonSpace += len(is)
						lines = append(lines, nextNonSpace)
						continue
					}
					nextNonSpace++
				}
				ch <- updateTokenLine(&Token{
					Value: text[spaceStart:nextNonSpace],
					Index: spaceStart,
					Text:  text[spaceStart:nextNonSpace],
					Kind:  spaceToken,
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
							Value: quoteValue,
							Index: quoteStart,
							Text:  text[quoteStart : j+1],
							Kind:  quoteToken,
						})
						quoteStart = -1
						quoteValue = ""
						break
					}

					quoteValue += string(c)
					j++
				}
				if j == len(text) {
					errCh <- errors.New("unclosed quote")
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

	return ch, errCh
}

var tokenKindNames = map[tokenKind]string{
	basicToken:       "basic",
	parensCloseToken: "close",
	parensOpenToken:  "open",
	quoteToken:       "quote",
	commentToken:     "comment",
	spaceToken:       "space",
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
