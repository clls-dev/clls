package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"github.com/clls-dev/clls/pkg/clls"
	"github.com/clls-dev/clls/pkg/lsp"
	"github.com/clls-dev/clls/pkg/lsph"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func main() {
	l := newLogger()
	l.Info("Logger initialized")

	docs := make(map[string]string)

	err := func() error {
		for {
			// Read header
			h, err := readHeader(l, os.Stdin)
			if err != nil {
				return errors.Wrap(err, "read header")
			}
			l.Debug("got lsp header", zap.Any("header", h))
			if h.ContentLength == nil {
				return errors.New("no Content-Length")
			}

			// Read message
			b := make([]byte, *h.ContentLength)
			if _, err := io.ReadFull(os.Stdin, b); err != nil {
				return errors.Wrap(err, "read message")
			}
			// l.Debug("got message", zap.String("content", string(b)))

			// Unmarshal request
			var req lsp.RequestMessage
			if err := json.Unmarshal(b, &req); err != nil {
				return errors.Wrap(err, "unmarshal message")
			}
			l.Debug("unmarshaled "+req.Method, zap.Any("request", req))

			params := req.Params.(map[string]interface{})

			// Respond
			switch req.Method {
			case "initialize":
				reply := lsp.ResponseMessage{
					Message: lsp.Message{Version: "2.0"},
					ID:      req.ID,
					Result: lsp.InitializeResult{
						ServerInfo: &lsp.InitializeResultServerInfo{
							Name:    "clls",
							Version: "0.1.0",
						},
						Capabilities: &lsp.ServerCapabilities{
							TextDocumentSync: lsp.Full,
							SemanticTokensProvider: lsp.SemanticTokensOptions{
								Legend: lsp.DefaultSemanticTokensLegend,
								Full:   true,
							},
						},
					},
				}
				if err := doReply(l, os.Stdout, &reply); err != nil {
					return errors.Wrap(err, "reply to initialize")
				}
				l.Debug("responded to initialize", zap.Any("reply", reply))
			case "textDocument/didOpen":
				td := params["textDocument"].(map[string]interface{})
				uriStr := td["uri"].(string)
				contentStr := td["text"].(string)
				docs[uriStr] = contentStr
			case "textDocument/didChange":
				td := params["textDocument"].(map[string]interface{})
				uriStr := td["uri"].(string)
				contentStr := params["contentChanges"].([]interface{})[0].(map[string]interface{})["text"].(string)
				docs[uriStr] = contentStr
			case "textDocument/didClose":
				tdi := params["textDocument"].(map[string]interface{})
				uriStr := tdi["uri"].(string)
				delete(docs, uriStr)
			case "textDocument/semanticTokens/full":
				const fileURIPrefix = "file://"
				uriStr := req.Params.(map[string]interface{})["textDocument"].(map[string]interface{})["uri"].(string)
				l.Debug("textDocument/semanticTokens/full " + uriStr)
				if !strings.HasPrefix(uriStr, fileURIPrefix) {
					if err := replyWithError(l, os.Stdout, req.ID, errors.New("textDocument.uri is not a file uri")); err != nil {
						return errors.Wrap(err, "reply with error to textDocument/semanticTokens/full")
					}
					continue
				}
				pathStr := strings.TrimPrefix(uriStr, fileURIPrefix)
				l.Debug("parsed uri", zap.String("path", pathStr))
				dirname := filepath.Dir(pathStr)
				filename := filepath.Base(pathStr)
				mod, err := clls.LoadCLVM(l, filename, func(p string) (string, error) {
					if d, ok := docs["file://"+filepath.Join(dirname, p)]; ok {
						return d, nil
					}
					b, err := ioutil.ReadFile(filepath.Join(dirname, p))
					if err != nil {
						return "", err
					}
					return string(b), err
				})
				if err != nil {
					l.Error("load clvm", zap.Error(err))
					reply := lsp.ResponseMessage{
						Message: lsp.Message{Version: "2.0"},
						ID:      req.ID,
						Result:  lsp.SemanticTokens{},
					}
					if err := doReply(l, os.Stdout, &reply); err != nil {
						return errors.Wrap(err, "reply to semantic tokens")
					}
					continue
				}
				l.Debug("parsed mod", zap.Any("val", mod))

				inserts := []insert(nil)

				if mod.IsMod {
					inserts = append(inserts, insert{Kind: "keyword", Token: mod.ModToken})
					if mod.Args != nil {
						inserts = insertParamsTokens(inserts, mod.Args)
					}
				}

				if mod.Constants != nil {
					for _, c := range mod.Constants {
						if c.Token != nil {
							inserts = append(inserts, insert{Kind: "keyword", Token: c.Token})
						}
						if t, ok := c.Name.(*clls.Token); ok && t != nil {
							inserts = append(inserts, insert{Kind: "variable", Modifiers: []string{"readonly"}, Token: t})
						}
						inserts = insertBody(inserts, c.Value)
					}
				}

				for _, t := range mod.Comments {
					inserts = append(inserts, insert{Kind: "comment", Token: t})
				}

				for _, incl := range mod.Includes {
					inserts = append(inserts, insert{Kind: "keyword", Token: incl.Token})
					if t, ok := incl.Value.(*clls.Token); ok && t != nil {
						inserts = append(inserts, insert{Kind: "string", Token: t})
					}
				}

				if len(mod.Functions) > 0 {
					for _, f := range mod.Functions {
						inserts = append(inserts, insert{Kind: "keyword", Token: f.KeywordToken})
						if f.Name != nil {
							inserts = append(inserts, insert{Kind: "function", Token: f.Name})
						}
						if f.Params != nil {
							inserts = insertParamsTokens(inserts, f.Params)
						}
						inserts = insertBody(inserts, f.Body)
					}
				}

				if mod.Main != nil {
					inserts = insertBody(inserts, mod.Main)
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

					l.Debug("generated tokens", zap.Any("inserts", inserts))

					ltoks := lsph.SemanticTokenSlice{}

					t := inserts[0].Token
					if t != nil {
						tt, tm, err := tokenInfo(l, inserts[0].Kind, inserts[0].Modifiers, &lsp.DefaultSemanticTokensLegend)
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

							tt, tm, err := tokenInfo(l, in.Kind, in.Modifiers, &lsp.DefaultSemanticTokensLegend)
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

				reply := lsp.ResponseMessage{
					Message: lsp.Message{Version: "2.0"},
					ID:      req.ID,
					Result:  lsp.SemanticTokens{Data: data},
				}
				if err := doReply(l, os.Stdout, &reply); err != nil {
					return errors.Wrap(err, "reply to semantic tokens")
				}
			case "shutdown":
				for k := range docs {
					delete(docs, k)
				}
				reply := lsp.ResponseMessage{
					Message: lsp.Message{Version: "2.0"},
					ID:      req.ID,
					Result:  nil,
				}
				if err := doReply(l, os.Stdout, &reply); err != nil {
					return errors.Wrap(err, "reply to shutdown")
				}
				l.Debug("responded to shutdown", zap.Any("reply", reply))
				return nil
			}
		}
	}()

	if err != nil {
		l.Error("main", zap.Error(err))
		panic(err)
	}

	l.Debug("done")
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
	case *clls.Token:
		if a.Value != "." {
			inserts = append(inserts, insert{Kind: "parameter", Modifiers: []string{"readonly"}, Token: a})
		}
	case *clls.ASTNode:
		for _, ac := range a.Children {
			inserts = insertParamsTokens(inserts, ac)
		}
	}
	return inserts
}

func insertBody(inserts []insert, node *clls.CodeBody) []insert {
	if node == nil {
		return inserts
	}
	switch node.Kind {
	case clls.IfBodyKind:
		if node.Token != nil {
			inserts = append(inserts, insert{Kind: "keyword", Token: node.Token})
		}
		inserts = insertBody(inserts, node.IfCond)
		inserts = insertBody(inserts, node.IfBranch)
		inserts = insertBody(inserts, node.ElseBranch)
	case clls.CallBodyKind:
		kind := "function"
		mods := []string(nil)
		if node.Function.Builtin {
			switch node.Function.Name.Value {
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
			inserts = insertBody(inserts, a)
		}
	case clls.OperatorBodyKind:
		if node.Token != nil {
			inserts = append(inserts, insert{Kind: "operator", Token: node.Token})
		}
		for _, child := range node.Children {
			inserts = insertBody(inserts, child)
		}
	case clls.ConstBodyKind:
		if node.Token != nil {
			inserts = append(inserts, insert{Kind: "variable", Modifiers: []string{"readonly"}, Token: node.Token})
		}
	case clls.VarBodyKind:
		if node.Token != nil {
			inserts = append(inserts, insert{Kind: "parameter", Token: node.Token})
		}
	case clls.FuncVarBodyKind:
		k := "function"
		mods := []string(nil)
		if node.Function.Builtin {
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
			inserts = insertBody(inserts, child)
		}
	}
	return inserts
}

func newLogger() *zap.Logger {
	cfg := zap.NewDevelopmentConfig()
	cfg.OutputPaths = []string{"/tmp/vqscode-clls/server.log"}
	l, err := cfg.Build()
	if err != nil {
		fmt.Println("failed to init logger:", err)
		l = zap.NewNop()
	}
	return l
}

type Header struct {
	ContentLength *int
	ContentType   *string
}

func (h *Header) write(l *zap.Logger, w io.Writer) error {
	kv := map[string]string{}
	if h.ContentLength != nil {
		kv["Content-Length"] = strconv.Itoa(*h.ContentLength)
	}
	if h.ContentType != nil {
		kv["Content-Type"] = *h.ContentType
	}

	lines := []string(nil)
	for k, v := range kv {
		lines = append(lines, k+": "+v)
	}

	hdr := strings.Join(lines, "\r\n") + "\r\n\r\n"

	l.Debug("marshaled header", zap.String("value", hdr))

	_, err := w.Write([]byte(hdr))
	return err
}

func readHeader(l *zap.Logger, r io.Reader) (*Header, error) {
	b := make([]byte, 1)
	ab := []byte(nil)
	for len(ab) < 4 || string(ab[len(ab)-4:]) != "\r\n\r\n" {
		_, err := io.ReadFull(r, b)
		if err != nil {
			return nil, errors.Wrap(err, "read until sep")
		}
		ab = append(ab, b[0])
	}
	str := string(ab[:len(ab)-4])

	lines := strings.Split(str, "\r\n")
	h := Header{}
	for _, line := range lines {
		assignIndex := strings.Index(line, ":")
		if assignIndex == -1 {
			return nil, fmt.Errorf("no separator in header line '%s'", line)
		}
		if assignIndex == len(line)-1 {
			return nil, fmt.Errorf("no value in header line '%s'", line)
		}
		if strings.HasPrefix(line, "Content-Length") {
			i, err := strconv.Atoi(strings.TrimSpace(line[assignIndex+1:]))
			if err != nil {
				return nil, errors.Wrap(err, fmt.Sprintf("parse Content-Length from '%s'", strings.TrimSpace(line[assignIndex+1:])))
			}
			h.ContentLength = &i
		} else if strings.HasPrefix(line, "Content-Type") {
			s := strings.TrimSpace(line[assignIndex+1:])
			h.ContentType = &s
		}
	}

	return &h, nil
}

func doReply(l *zap.Logger, w io.Writer, res *lsp.ResponseMessage) error {
	replyBytes, err := json.Marshal(res)
	if err != nil {
		return errors.Wrap(err, "marshal reply")
	}
	l.Debug("marshaled reply", zap.String("value", string(replyBytes)))
	cl := len(replyBytes)
	if err := (&Header{ContentLength: &cl}).write(l, os.Stdout); err != nil {
		return errors.Wrap(err, "write header")
	}
	if _, err := os.Stdout.Write(replyBytes); err != nil {
		return errors.Wrap(err, "write message")
	}
	return nil
}

func replyWithError(l *zap.Logger, w io.Writer, id lsp.IntegerOrString, err error) error {
	l.Error(fmt.Sprintf("replying to %s with error", id), zap.Error(err))
	return doReply(l, w, &lsp.ResponseMessage{
		Message: lsp.Message{Version: "2.0"},
		ID:      id,
		Error: &lsp.ResponseError{
			Code:    lsp.UnknownErrorCode,
			Message: err.Error(),
		},
	})
}

type insert struct {
	Kind      string
	Modifiers []string
	Token     *clls.Token
}
