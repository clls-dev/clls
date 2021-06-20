package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/clls-dev/clls/pkg/lspsrv"
	"github.com/pkg/errors"
	lsp "go.lsp.dev/protocol"
	"go.uber.org/zap"
)

func newLogger() *zap.Logger {
	cfg := zap.NewDevelopmentConfig()
	var err error
	var l *zap.Logger
	if err = os.MkdirAll("/tmp/vscode-clls", os.ModePerm); err == nil {
		cfg.OutputPaths = []string{"/tmp/vscode-clls/server.log"}
		l, err = cfg.Build()
	}
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
				if err := lspsrv.ReplyWithErrorCode(l, out, req.ID, errors.New("shutdown mode"), lspsrv.InvalidRequest); err != nil {
					return errors.Wrap(err, "reply to request in shutdown mode")
				}
				continue
			}

			// Unmarshal params
			params, err := lspsrv.Unmarshal(req.Method, req.Params)
			if err != nil {
				if req.Method == "initialize" {
					params = &lsp.InitializeParams{}
					l.Error("unmarshal initialize params", zap.Error(err))
				} else {
					if err := lspsrv.ReplyWithError(l, out, req.ID, errors.Wrap(err, fmt.Sprintf("unmarshal '%s'", req.Method))); err != nil {
						return errors.Wrap(err, "reply with unmarshal error")
					}
					continue
				}
			}

			// Handle
			ctx, cancel := context.WithCancel(context.TODO())
			reply, err := srv.Request(ctx, req.Method, params)
			if err != nil {
				cancel()
				if err := lspsrv.ReplyWithError(l, out, req.ID, errors.Wrap(err, fmt.Sprintf("handle '%s'", req.Method))); err != nil {
					return errors.Wrap(err, "reply with handle error")
				}
				continue
			}
			cancel()

			// Maybe reply
			if req.ID != nil && reply != nil {
				if err := lspsrv.Reply(l, out, &lspsrv.ResponseMessage{
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
