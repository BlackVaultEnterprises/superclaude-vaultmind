# SuperSoldierAI (SuperClaude) Integration Guide

## Current Status

We've built the SuperClaude enhancement layer for OpenCode, but it requires manual integration into the OpenCode TUI to work exactly like Claude Code + SuperClaude.

## What We Have

✅ **Complete SuperClaude Implementation**
- All 18 commands with templates
- 9 cognitive personas
- Universal flag system
- Command parsing and routing
- Tests passing (after fixes)

✅ **Configuration**
- Ollama integration configured
- Local model support ready

## What's Missing for Full Integration

The integration point in `internal/tui/page/chat.go` needs to be manually added. The patch file shows what needs to be done:

1. Import the superclaude package
2. In the `sendMessage` function, check if input is a SuperClaude command
3. If yes, handle it through SuperClaude; if no, process normally

## To Make It Work EXACTLY Like Claude Code + SuperClaude

### Option 1: Manual Integration (Recommended)
```bash
# 1. Edit internal/tui/page/chat.go
# 2. Add import: "github.com/opencode-ai/opencode/internal/superclaude"
# 3. In sendMessage function, add the SuperClaude handler check
# 4. Build: go build
# 5. Run: ./opencode
```

### Option 2: Fork Approach
```bash
# 1. Fork OpenCode
# 2. Add SuperClaude as a proper module
# 3. Integrate at the command processing level
# 4. Submit PR to OpenCode (if they're interested)
```

### Option 3: Wrapper Approach
Create a wrapper that intercepts commands before sending to OpenCode:
```go
// supersoldier-wrapper.go
// Intercept user input
// If SuperClaude command → process with our handler
// Else → pass through to OpenCode
```

## The Reality

OpenCode wasn't designed with a plugin system for command extensions. To get **EXACT** Claude Code + SuperClaude behavior, you need to:

1. Modify OpenCode's source (as shown in the patch)
2. Maintain your own fork
3. Or convince OpenCode to add a plugin/extension system

## Quick Test Without Full Integration

You can test the SuperClaude command processing:
```go
// test-superclaude.go
package main

import (
    "fmt"
    "github.com/opencode-ai/opencode/internal/superclaude"
)

func main() {
    input := "/user:build --react --tdd"
    parsed, _ := superclaude.ParseSuperClaudeCommand(input)
    fmt.Printf("Command: %s, Persona: %s\n", parsed.Command, parsed.Flags.Persona)
}
```

## Bottom Line

The SuperClaude layer is fully implemented and tested. To work exactly like Claude Code + SuperClaude, it needs to be integrated into OpenCode's message handling loop. The patch file shows exactly where and how, but it requires modifying OpenCode's source code.