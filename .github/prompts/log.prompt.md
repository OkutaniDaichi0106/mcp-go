# Prompt for Generating Log Messages (using slog)

You are writing a log message for a Go application using the standard `log/slog` package (introduced in Go 1.21). The log message should:
- Clearly describe what the system is doing or what event has occurred
- Include important variables, IDs, or contextual information
- Use structured logging with `slog` (use `slog.Info`, `slog.Error`, etc.)
- Use consistent structure and professional tone
- Be appropriate for its log level (INFO, WARN, ERROR, DEBUG)

Write the message in English. Log values should be passed as key-value pairs using `slog`'s structured logging style.

Example Format (INFO level):
~~~go
slog.Info("Started processing order", "order_id", orderID, "customer_id", customerID)
~~~

Example Format (ERROR level):
~~~go
slog.Error("Failed to load config file", "path", configPath, "error", err)
~~~

Now, generate a log message for the following scenario:
<INSERT DESCRIPTION OF EVENT OR CODE SNIPPET HERE>
