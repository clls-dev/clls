package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"unicode"

	"github.com/clls-dev/clls/pkg/clls"
	"github.com/clls-dev/clls/pkg/lsp"
	"github.com/clls-dev/clls/pkg/lspsrv"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const fileURIPrefix = "file://"

func main() {
	l := newLogger()
	l.Info("Logger initialized")

	err := func() error {
		srv := newServer(l)
		for !srv.exit {
			// Read header
			h, err := readHeader(l, os.Stdin)
			if err != nil {
				if err == io.EOF {
					break
				}
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
			var req lsp.RawRequestMessage
			if err := json.Unmarshal(b, &req); err != nil {
				return errors.Wrap(err, "unmarshal message")
			}
			l.Debug("unmarshaled raw message", zap.String("method", req.Method))

			// Respond
			reply, err := lspsrv.LanguageServerHandle(srv, req.Method, req.Params)
			if err != nil {
				if err := replyWithError(l, os.Stdout, req.ID, errors.Wrap(err, fmt.Sprintf("handle '%s'", req.Method))); err != nil {
					return errors.Wrap(err, "reply with error")
				}
			}
			if req.ID != nil && reply != nil {
				if err := doReply(l, os.Stdout, &lsp.ResponseMessage{
					ID:     req.ID,
					Result: reply,
				}); err != nil {
					return errors.Wrap(err, "reply")
				}
			}
		}

		return nil
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

func insertBody(inserts []insert, node *clls.CodeBody, funcsByName map[string]*clls.Function) []insert {
	if node == nil {
		return inserts
	}
	switch node.Kind {
	case clls.IfBodyKind:
		if node.Token != nil {
			inserts = append(inserts, insert{Kind: "keyword", Token: node.Token})
		}
		inserts = insertBody(inserts, node.IfCond, funcsByName)
		inserts = insertBody(inserts, node.IfBranch, funcsByName)
		inserts = insertBody(inserts, node.ElseBranch, funcsByName)
	case clls.CallBodyKind:
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
	case clls.OperatorBodyKind:
		if node.Token != nil {
			inserts = append(inserts, insert{Kind: "operator", Token: node.Token})
		}
		for _, child := range node.Children {
			inserts = insertBody(inserts, child, funcsByName)
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
			if err == io.EOF {
				return nil, err
			}
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
	if res.Message.Version == "" {
		res.Message.Version = "2.0"
	}
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
	l.Error(fmt.Sprintf("replying to %v with error", id), zap.Error(err))
	return doReply(l, w, &lsp.ResponseMessage{
		ID: id,
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
