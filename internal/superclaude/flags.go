package superclaude

import (
	"fmt"
	"strings"
)

// Flags represents all possible SuperClaude flags
type Flags struct {
	// Core flags
	Persona         string // --persona-architect
	Think           string // --think, --think-hard, --ultrathink
	UltraCompressed bool   // --uc
	Plan            bool   // --plan
	Evidence        bool   // --evidence (--c7 equivalent)
	ValidationOnly  bool   // --validate
	Sequential      bool   // --seq
	AllMCP          bool   // --all-mcp

	// Additional dynamic flags
	Additional map[string]string // For command-specific flags
}

// ParsedCommand represents a fully parsed SuperClaude command
type ParsedCommand struct {
	Command  string
	Target   string
	Flags    *Flags
	RawInput string
}

// ParseSuperClaudeCommand parses a SuperClaude command string
func ParseSuperClaudeCommand(input string) (*ParsedCommand, error) {
	input = strings.TrimSpace(input)

	// Check if it's a SuperClaude command
	if !strings.HasPrefix(input, "/user:") && !strings.HasPrefix(input, "/persona:") {
		return nil, fmt.Errorf("not a SuperClaude command")
	}

	// Extract persona if specified with /persona:
	var persona string
	if strings.HasPrefix(input, "/persona:") {
		parts := strings.SplitN(input, " ", 2)
		if len(parts) < 2 {
			return nil, fmt.Errorf("invalid persona command format")
		}
		personaPart := strings.TrimPrefix(parts[0], "/persona:")
		persona = personaPart
		input = parts[1] // Continue with the rest

		// Check for arrow operator
		if strings.HasPrefix(input, "→") || strings.HasPrefix(input, "->") {
			input = strings.TrimSpace(strings.TrimPrefix(strings.TrimPrefix(input, "→"), "->"))
		}
	}

	// Now parse the /user: command
	if !strings.HasPrefix(input, "/user:") {
		return nil, fmt.Errorf("expected /user: command after persona")
	}

	// Split the command from flags
	parts := strings.Fields(input)
	if len(parts) < 1 {
		return nil, fmt.Errorf("empty command")
	}

	// Extract command and target
	commandPart := strings.TrimPrefix(parts[0], "/user:")
	var command, target string

	// The command is the first part after /user:
	command = commandPart

	// Parse flags and extract target
	flags := &Flags{
		Persona:    persona,
		Additional: make(map[string]string),
	}

	var targetParts []string
	i := 1

	for i < len(parts) {
		part := parts[i]

		if strings.HasPrefix(part, "--") {
			// Parse flag
			flagName := strings.TrimPrefix(part, "--")

			switch flagName {
			// Thinking modes
			case "think":
				flags.Think = "standard"
			case "think-hard":
				flags.Think = "deep"
			case "ultrathink":
				flags.Think = "ultra"

			// Boolean flags
			case "uc", "ultracompressed":
				flags.UltraCompressed = true
			case "plan":
				flags.Plan = true
			case "evidence", "c7":
				flags.Evidence = true
			case "validate", "validation-only":
				flags.ValidationOnly = true
			case "seq", "sequential":
				flags.Sequential = true
			case "all-mcp":
				flags.AllMCP = true

			default:
				// Check for persona flags
				if strings.HasPrefix(flagName, "persona-") {
					flags.Persona = strings.TrimPrefix(flagName, "persona-")
				} else {
					// Check if next part is a value
					if i+1 < len(parts) && !strings.HasPrefix(parts[i+1], "--") {
						i++
						flags.Additional[flagName] = parts[i]
					} else {
						flags.Additional[flagName] = "true"
					}
				}
			}
		} else {
			// It's part of the target
			targetParts = append(targetParts, part)
		}
		i++
	}

	target = strings.Join(targetParts, " ")

	// If no persona specified, use the default for the command
	if flags.Persona == "" {
		// Check if specific flags should influence persona selection
		if command == "build" && flags.Additional["react"] == "true" {
			flags.Persona = "frontend"
		} else {
			flags.Persona = GetPersonaForCommand(command)
		}
	}

	return &ParsedCommand{
		Command:  command,
		Target:   target,
		Flags:    flags,
		RawInput: input,
	}, nil
}

// GetThinkingTokens returns the max tokens based on thinking mode
func GetThinkingTokens(thinkMode string) int {
	switch thinkMode {
	case "ultra":
		return 32000
	case "deep":
		return 16000
	case "standard":
		return 8000
	default:
		return 4000
	}
}

// ValidateFlags ensures flags are compatible
func (f *Flags) Validate() error {
	// Validate persona
	if f.Persona != "" {
		if _, exists := Personas[f.Persona]; !exists {
			return fmt.Errorf("unknown persona: %s", f.Persona)
		}
	}

	// Validate thinking mode
	if f.Think != "" && f.Think != "standard" && f.Think != "deep" && f.Think != "ultra" {
		return fmt.Errorf("invalid thinking mode: %s", f.Think)
	}

	// Check for conflicting flags
	if f.ValidationOnly && f.Plan {
		return fmt.Errorf("cannot use --validate and --plan together")
	}

	return nil
}

// MergeFlags combines two flag sets with precedence to the second
func MergeFlags(base, override *Flags) *Flags {
	result := &Flags{
		Persona:         override.Persona,
		Think:           override.Think,
		UltraCompressed: override.UltraCompressed,
		Plan:            override.Plan,
		Evidence:        override.Evidence,
		ValidationOnly:  override.ValidationOnly,
		Sequential:      override.Sequential,
		AllMCP:          override.AllMCP,
		Additional:      make(map[string]string),
	}

	// If override doesn't specify, use base
	if result.Persona == "" {
		result.Persona = base.Persona
	}
	if result.Think == "" {
		result.Think = base.Think
	}
	if !result.UltraCompressed {
		result.UltraCompressed = base.UltraCompressed
	}
	if !result.Plan {
		result.Plan = base.Plan
	}
	if !result.Evidence {
		result.Evidence = base.Evidence
	}

	// Merge additional flags
	for k, v := range base.Additional {
		result.Additional[k] = v
	}
	for k, v := range override.Additional {
		result.Additional[k] = v
	}

	return result
}
