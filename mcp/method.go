package mcp

type Method string

const (
	MethodInit Method = "initialize"

	// Tools
	MethodListTools         Method = "tools/list"
	MethodCallTool          Method = "tools/call"
	MethodNotifyToolChanged Method = "notifications/tools/list_changed"

	// Resources
	MethodListResources         Method = "resources/list"
	MethodReadResource          Method = "resources/read"
	MethodSubscribeResource     Method = "resources/subscribe"
	MethodNotifyResourceChanged Method = "notifications/resources/list_changed"
	MethodNotifyResourceUpdated Method = "notifications/resources/updated"

	// Prompts
	MethodListPrompts         Method = "prompts/list"
	MethodGetPrompt           Method = "prompts/get"
	MethodNotifyPromptChanged Method = "notifications/prompts/list_changed"

	// Logging
	MethodSetLogLevel Method = "logging/setLevel"

	// Sampling
	MethodCreateSampleMessage Method = "sampling/createMessage"

	// Roots
	MethodListRoots         Method = "roots/list"
	MethodNotifyRootChanged Method = "notifications/roots/list_changed"
)
