# pkg/messages/ - Rule Message Templates

## Purpose
Centralized message definitions for all KTN rules. Provides short and verbose formats.

## Structure
```
messages/
├── messages.go   # Registry + GetMessage() function
├── func.go       # KTN-FUNC-* messages
├── var.go        # KTN-VAR-* messages
├── const.go      # KTN-CONST-* messages
├── struct.go     # KTN-STRUCT-* messages
├── comment.go    # KTN-COMMENT-* messages
├── test.go       # KTN-TEST-* messages
└── *_test.go     # Tests
```

## Message Format
```go
var Messages = map[string]Message{
    "KTN-FUNC-001": {
        Short:   "Function too long (%d lines, max %d)",
        Verbose: "Function '%s' has %d lines which exceeds the maximum of %d. Consider extracting helper functions.",
    },
}
```

## Usage
```go
msg := messages.Get("KTN-FUNC-001")
formatted := fmt.Sprintf(msg.Short, actualLines, maxLines)
```

## Conventions
- Short: One line, factual (used in lint output)
- Verbose: Explanation + suggestion (used with --verbose flag)
- Use format verbs (%s, %d) for dynamic values
