package mcp

type PromptWriter interface {
	Write(role Role, content Content) error
	CloseWithError(code ErrorCode, msg string) error
}

type Role string

const (
	User      Role = "user"
	Assistant Role = "assistant"
)
