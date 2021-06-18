package lspsrv

import (
	"encoding/json"
	"io"
	"os"

	"github.com/clls-dev/clls/pkg/lsp"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type FileTransport struct {
	in  *os.File
	out *os.File
	l   *zap.Logger
}

func NewFileTransport(l *zap.Logger, in *os.File, out *os.File) *FileTransport {
	ft := FileTransport{in, out, l}
	return &ft
}

func (ft *FileTransport) Recv() (*lsp.RawRequestMessage, error) {
	// Read header
	h, err := ReadHeader(ft.l, ft.in)
	if err != nil {
		return nil, err
	}
	if h.ContentLength == nil {
		return nil, errors.New("no Content-Length")
	}

	// Read message
	b := make([]byte, *h.ContentLength)
	if _, err := io.ReadFull(ft.in, b); err != nil {
		return nil, err
	}

	// Unmarshal request
	var req lsp.RawRequestMessage
	if err := json.Unmarshal(b, &req); err != nil {
		return nil, err
	}

	return &req, nil
}

func (ft *FileTransport) Send(res *lsp.ResponseMessage) error {
	return Reply(ft.l, ft.out, res)
}

func Reply(l *zap.Logger, w io.Writer, res *lsp.ResponseMessage) error {
	if res.Message.Version == "" {
		res.Message.Version = "2.0"
	}
	replyBytes, err := json.Marshal(res)
	if err != nil {
		return errors.Wrap(err, "marshal reply")
	}
	cl := len(replyBytes)
	if err := (&Header{ContentLength: &cl}).Write(l, w); err != nil {
		return errors.Wrap(err, "write header")
	}
	if _, err := w.Write(replyBytes); err != nil {
		return errors.Wrap(err, "write message")
	}
	return nil
}

func ReplyWithError(l *zap.Logger, w io.Writer, id lsp.IntegerOrString, err error) error {
	return ReplyWithErrorCode(l, w, id, err, lsp.UnknownErrorCode)
}

func ReplyWithErrorCode(l *zap.Logger, w io.Writer, id lsp.IntegerOrString, err error, code lsp.ErrorCode) error {
	l.Error("replying with error", zap.Error(err), zap.Any("code", code))
	return Reply(l, w, &lsp.ResponseMessage{
		ID: id,
		Error: &lsp.ResponseError{
			Code:    code,
			Message: err.Error(),
		},
	})
}
