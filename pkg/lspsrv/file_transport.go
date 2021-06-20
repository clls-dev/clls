package lspsrv

import (
	"encoding/json"
	"io"
	"os"

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

type RawRequestMessage struct {
	Version string          `json:"jsonrpc"`
	ID      interface{}     `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
}

type ResponseError struct {
	Code    ErrorCode   `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ResponseMessage struct {
	Version string         `json:"jsonrpc"`
	ID      interface{}    `json:"id,omitempty"`
	Result  interface{}    `json:"result,omitempty"`
	Error   *ResponseError `json:"error,omitempty"`
}

func (ft *FileTransport) Recv() (*RawRequestMessage, error) {
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
	var req RawRequestMessage
	if err := json.Unmarshal(b, &req); err != nil {
		return nil, err
	}

	return &req, nil
}

func (ft *FileTransport) Send(res *ResponseMessage) error {
	return Reply(ft.l, ft.out, res)
}

func Reply(l *zap.Logger, w io.Writer, res *ResponseMessage) error {
	if res.Version == "" {
		res.Version = "2.0"
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

func ReplyWithError(l *zap.Logger, w io.Writer, id interface{}, err error) error {
	return ReplyWithErrorCode(l, w, id, err, UnknownErrorCode)
}

func ReplyWithErrorCode(l *zap.Logger, w io.Writer, id interface{}, err error, code ErrorCode) error {
	l.Error("replying with error", zap.Error(err), zap.Any("code", code))
	return Reply(l, w, &ResponseMessage{
		ID: id,
		Error: &ResponseError{
			Code:    code,
			Message: err.Error(),
		},
	})
}

type ErrorCode int64

const (
	ParseError                     = ErrorCode(-32700)
	InvalidRequest                 = ErrorCode(-32600)
	MethodNotFound                 = ErrorCode(-32601)
	InvalidParams                  = ErrorCode(-32602)
	InternalError                  = ErrorCode(-32603)
	jsonrpcReservedErrorRangeStart = ErrorCode(-32099)
	/** @deprecated use jsonrpcReservedErrorRangeStart */
	// serverErrorStart = ErrorCode(jsonrpcReservedErrorRangeStart)

	/**
	 * Error code indicating that a server received a notification or
	 * request before the server has received the `initialize` request.
	 */
	ServerNotInitialized = ErrorCode(-32002)
	UnknownErrorCode     = ErrorCode(-32001)

	/**
	 * This is the start range of JSON RPC reserved error codes.
	 * It doesn't denote a real error code.
	 *
	 * @since 3.16.0
	 */
	// jsonrpcReservedErrorRangeEnd = ErrorCode(-32000)
	/** @deprecated use jsonrpcReservedErrorRangeEnd */
	// serverErrorEnd = ErrorCode(jsonrpcReservedErrorRangeEnd)

	/**
	 * This is the start range of LSP reserved error codes.
	 * It doesn't denote a real error code.
	 *
	 * @since 3.16.0
	 */
	// lspReservedErrorRangeStart = ErrorCode(-32899)

	ContentModified  = ErrorCode(-32801)
	RequestCancelled = ErrorCode(-32800)

	/**
	 * This is the end range of LSP reserved error codes.
	 * It doesn't denote a real error code.
	 *
	 * @since 3.16.0
	 */
	// lspReservedErrorRangeEnd = ErrorCode(-32800)
)
