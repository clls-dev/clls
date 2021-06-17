package main

import (
	"fmt"
	"io"
	"os"

	"github.com/clls-dev/clls/pkg/lsp"
	"github.com/clls-dev/clls/pkg/lspsrv"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const fileURIPrefix = "file://"

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

func main() {
	l := newLogger()
	l.Info("Logger initialized")

	out := os.Stdout

	err := func() error {
		srv := newServer(l.Named("ls"))
		transport := lspsrv.NewFileTransport(l.Named("trs"), os.Stdin, out)
		l = l.Named("loop")
		for !srv.exit {
			// Read message
			req, err := transport.Recv()
			if err != nil {
				if err == io.EOF {
					break
				}
				return errors.Wrap(err, "recv")
			}

			l.Debug("recv", zap.String("method", req.Method))

			// Ignore requests in shutdown mode
			if srv.down && req.Method != "exit" {
				if req.ID == nil {
					continue
				}
				if err := lspsrv.ReplyWithErrorCode(l, out, req.ID, errors.New("shutdown mode"), lsp.InvalidRequest); err != nil {
					return errors.Wrap(err, "reply to request in shutdown mode")
				}
				continue
			}

			// Respond
			reply, err := lspsrv.LanguageServerHandle(srv, req.Method, req.Params)
			if err != nil {
				if err := lspsrv.ReplyWithError(l, out, req.ID, errors.Wrap(err, fmt.Sprintf("handle '%s'", req.Method))); err != nil {
					return errors.Wrap(err, "reply with error")
				}
			}
			if req.ID != nil && reply != nil {
				if err := lspsrv.Reply(l, out, &lsp.ResponseMessage{
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
