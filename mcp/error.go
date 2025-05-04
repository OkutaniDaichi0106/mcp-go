package mcp

var _ error = (*Error)(nil)

type Error struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Data    any       `json:"data,omitempty"`
}

func (e *Error) Error() string {
	return e.Message
}

type ErrorCode int

const (
	ParseErrorCode           ErrorCode = -32700
	InvalidRequestErrorCode  ErrorCode = -32602
	MethodNotFoundErrorCode  ErrorCode = -32601
	InvalidParamsErrorCode   ErrorCode = -32602
	JSONRPCInternalErrorCode ErrorCode = -32603

	MCPInternalErrorCode ErrorCode = -32000

	ErrToolNotFoundCode     ErrorCode = MCPInternalErrorCode - 1 // -32001
	ErrResourceNotFoundCode ErrorCode = MCPInternalErrorCode - 2 // -32002
	ErrPromptNotFoundCode   ErrorCode = MCPInternalErrorCode - 3 // -32003
	ErrRootNotFoundCode     ErrorCode = MCPInternalErrorCode - 4 // -32004
	ErrSampleNotFoundCode   ErrorCode = MCPInternalErrorCode - 5 // -32005
)

var (
	ErrJSONRPCInternalError = &Error{
		Code:    JSONRPCInternalErrorCode,
		Message: "Internal JSON-RPC error",
	}
	ErrInvalidParams = &Error{
		Code:    InvalidParamsErrorCode,
		Message: "Invalid params",
	}
	ErrMethodNotFound = &Error{
		Code:    MethodNotFoundErrorCode,
		Message: "Method not found",
	}
	ErrInvalidRequest = &Error{
		Code:    InvalidRequestErrorCode,
		Message: "Invalid request",
	}
	ErrParseError = &Error{
		Code:    ParseErrorCode,
		Message: "Parse error",
	}

	ErrInternalError = &Error{
		Code:    MCPInternalErrorCode,
		Message: "Internal MCP error",
	}

	ErrToolNotFound = &Error{
		Code:    ErrToolNotFoundCode,
		Message: "Tool not found",
	}

	ErrResourceNotFound = &Error{
		Code:    ErrResourceNotFoundCode,
		Message: "Resource not found",
	}

	ErrPromptNotFound = &Error{
		Code:    ErrPromptNotFoundCode,
		Message: "Prompt not found",
	}

	ErrRootNotFound = &Error{
		Code:    ErrRootNotFoundCode,
		Message: "Root not found",
	}

	ErrSampleNotFound = &Error{
		Code:    ErrSampleNotFoundCode,
		Message: "Sample not found",
	}
)
