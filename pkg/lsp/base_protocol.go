package lsp

import "encoding/json"

type Integer = int64
type UInteger = uint32
type Decimal = float64

type IntegerOrString interface{}
type IntegerOrNull interface{}
type StringOrNull interface{}
type ArrayOrObject interface{}

type Message struct {
	Version string `json:"jsonrpc"`
}
type RequestMessage struct {
	Message
	ID     IntegerOrString `json:"id"`
	Method string          `json:"method"`
	Params ArrayOrObject   `json:"params"`
}
type RawRequestMessage struct {
	Message
	ID     IntegerOrString `json:"id"`
	Method string          `json:"method"`
	Params json.RawMessage `json:"params"`
}
type ResponseMessage struct {
	Message
	ID     IntegerOrString `json:"id"`
	Result interface{}     `json:"result,omitempty"`
	Error  *ResponseError  `json:"error,omitempty"`
}
type ResponseError struct {
	Code    ErrorCode   `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
type NotificationMessage struct {
	Method string        `json:"method"`
	Params ArrayOrObject `json:"params"`
}
type CancelParams struct {
	ID IntegerOrString `json:"id"`
}
type ProgressParams struct {
	Token ProgressToken `json:"token"`
	Value interface{}   `json:"value"`
}

type ProgressToken = IntegerOrString

type ErrorCode Integer

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
