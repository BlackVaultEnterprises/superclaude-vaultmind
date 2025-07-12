package superclaude

import (
	"context"
	"fmt"
	"strings"

	"github.com/opencode-ai/opencode/internal/llm/agent"
	"github.com/opencode-ai/opencode/internal/logging"
)

// GetAvailableCommands returns all available SuperClaude commands
func GetAvailableCommands() []string {
	commands := make([]string, 0, len(Commands))
	for cmd := range Commands {
		commands = append(commands, cmd)
	}
	return commands
}

// GetAvailablePersonas returns all available personas
func GetAvailablePersonas() []string {
	personas := make([]string, 0, len(Personas))
	for name := range Personas {
		personas = append(personas, name)
	}
	return personas
}

// GetAvailableFlags returns all available flags
func GetAvailableFlags() []string {
	return []string{
		"--uc", "--ultracompressed",
		"--think", "--think-hard", "--ultrathink",
		"--plan",
		"--evidence", "--c7",
		"--validate", "--validation-only",
		"--seq", "--sequential",
		"--all-mcp",
		"--persona-<name>",
	}
}

// SuperClaudeHandler handles SuperClaude commands within OpenCode
type SuperClaudeHandler struct {
	agent agent.Service
}

// NewSuperClaudeHandler creates a new SuperClaude handler
func NewSuperClaudeHandler(agent agent.Service) *SuperClaudeHandler {
	return &SuperClaudeHandler{
		agent: agent,
	}
}

// HandleCommand processes a potential SuperClaude command
func (h *SuperClaudeHandler) HandleCommand(ctx context.Context, sessionID string, input string) (bool, error) {
	// Validate inputs
	if sessionID == "" {
		return false, fmt.Errorf("session ID is required")
	}
	
	// Try to parse as SuperClaude command
	parsed, err := ParseSuperClaudeCommand(input)
	if err != nil {
		// Not a SuperClaude command, let OpenCode handle it normally
		return false, nil
	}

	// Validate flags
	if err := parsed.Flags.Validate(); err != nil {
		return true, fmt.Errorf("invalid flags: %w", err)
	}

	// Get the command
	cmd, exists := Commands[parsed.Command]
	if !exists {
		return true, fmt.Errorf("unknown command: %s", parsed.Command)
	}

	// Get the persona
	persona := GetPersona(parsed.Flags.Persona)

	// Build the enhanced prompt
	prompt, err := cmd.BuildPrompt(persona, parsed.Flags, parsed.Target, parsed.RawInput)
	if err != nil {
		return true, fmt.Errorf("failed to build prompt: %w", err)
	}

	// Apply thinking mode by adjusting context
	if parsed.Flags.Think != "" {
		prompt = applyThinkingMode(prompt, parsed.Flags.Think)
	}

	// Apply ultra-compressed mode
	if parsed.Flags.UltraCompressed {
		prompt = applyUltraCompressed(prompt)
	}

	// Log the SuperClaude command execution
	logging.Info("Executing SuperClaude command",
		"command", parsed.Command,
		"persona", persona.Name,
		"target", parsed.Target,
		"flags", formatFlags(parsed.Flags))

	// Execute through the agent with the enhanced prompt
	events, err := h.agent.Run(ctx, sessionID, prompt)
	if err != nil {
		return true, err
	}

	// Handle the response events
	go h.handleAgentEvents(events, parsed)

	return true, nil
}

// applyThinkingMode enhances the prompt for different thinking levels
func applyThinkingMode(prompt string, thinkMode string) string {
	prefix := ""

	switch thinkMode {
	case "ultra":
		prefix = `ULTRATHINK MODE ACTIVATED:
- Perform exhaustive analysis with multiple perspectives
- Consider edge cases, failure modes, and optimization opportunities
- Provide detailed reasoning chains for all decisions
- Challenge assumptions and explore alternatives

`
	case "deep":
		prefix = `DEEP THINKING MODE:
- Analyze thoroughly before implementation
- Consider architectural implications
- Evaluate trade-offs explicitly

`
	case "standard":
		prefix = `THINKING MODE:
- Plan before executing
- Consider key design decisions

`
	}

	return prefix + prompt
}

// applyUltraCompressed adds compression instructions
func applyUltraCompressed(prompt string) string {
	return prompt + `

ULTRA-COMPRESSED MODE:
- Use 70% fewer tokens
- No explanations unless critical
- Code > commentary
- Terse, efficient output`
}

// formatFlags formats flags for logging
func formatFlags(flags *Flags) string {
	var parts []string

	if flags.Persona != "" {
		parts = append(parts, fmt.Sprintf("persona=%s", flags.Persona))
	}
	if flags.Think != "" {
		parts = append(parts, fmt.Sprintf("think=%s", flags.Think))
	}
	if flags.UltraCompressed {
		parts = append(parts, "ultracompressed")
	}
	if flags.Evidence {
		parts = append(parts, "evidence")
	}
	if flags.Sequential {
		parts = append(parts, "sequential")
	}

	for k, v := range flags.Additional {
		parts = append(parts, fmt.Sprintf("%s=%s", k, v))
	}

	return strings.Join(parts, ", ")
}

// handleAgentEvents processes events from the agent
func (h *SuperClaudeHandler) handleAgentEvents(events <-chan agent.AgentEvent, parsed *ParsedCommand) {
	for event := range events {
		switch event.Type {
		case agent.AgentEventTypeResponse:
			// Response handled by OpenCode's UI
			continue
		case agent.AgentEventTypeError:
			logging.Error("SuperClaude command error",
				"command", parsed.Command,
				"error", event.Error)
		}
	}
}

// GetMaxTokensForCommand returns appropriate token limit based on command and flags
func GetMaxTokensForCommand(parsed *ParsedCommand) int {
	baseTokens := 4096

	// Adjust based on thinking mode
	if parsed.Flags.Think != "" {
		baseTokens = GetThinkingTokens(parsed.Flags.Think)
	}

	// Adjust based on command type
	switch parsed.Command {
	case "analyze", "design", "explain":
		baseTokens = int(float64(baseTokens) * 1.5)
	case "build", "improve", "refactor":
		baseTokens = int(float64(baseTokens) * 2)
	}

	// Reduce for ultra-compressed mode
	if parsed.Flags.UltraCompressed {
		baseTokens = int(float64(baseTokens) * 0.7)
	}

	return baseTokens
}
