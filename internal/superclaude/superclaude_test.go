package superclaude

import (
	"testing"
)

func TestParseSuperClaudeCommand(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantCmd     string
		wantPersona string
		wantTarget  string
		wantErr     bool
	}{
		{
			name:        "Basic command",
			input:       "/user:build --react",
			wantCmd:     "build",
			wantPersona: "frontend", // Default for build with --react
			wantTarget:  "",
			wantErr:     false,
		},
		{
			name:        "Command with target",
			input:       "/user:analyze codebase --deep",
			wantCmd:     "analyze",
			wantPersona: "analyzer", // Default for analyze
			wantTarget:  "codebase",
			wantErr:     false,
		},
		{
			name:        "Persona override",
			input:       "/persona:security → /user:scan --owasp",
			wantCmd:     "scan",
			wantPersona: "security",
			wantTarget:  "",
			wantErr:     false,
		},
		{
			name:        "Multiple flags",
			input:       "/user:test --coverage --e2e --pup",
			wantCmd:     "test",
			wantPersona: "qa",
			wantTarget:  "",
			wantErr:     false,
		},
		{
			name:        "Thinking mode",
			input:       "/user:design api --ddd --think-hard",
			wantCmd:     "design",
			wantPersona: "architect",
			wantTarget:  "api",
			wantErr:     false,
		},
		{
			name:        "Not a SuperClaude command",
			input:       "regular message",
			wantCmd:     "",
			wantPersona: "",
			wantTarget:  "",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsed, err := ParseSuperClaudeCommand(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("ParseSuperClaudeCommand() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil {
				return
			}

			if parsed.Command != tt.wantCmd {
				t.Errorf("Command = %v, want %v", parsed.Command, tt.wantCmd)
			}

			if parsed.Target != tt.wantTarget {
				t.Errorf("Target = %v, want %v", parsed.Target, tt.wantTarget)
			}

			if parsed.Flags.Persona != tt.wantPersona {
				t.Errorf("Persona = %v, want %v", parsed.Flags.Persona, tt.wantPersona)
			}
		})
	}
}

func TestFlags(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		checkFlags func(*Flags) bool
	}{
		{
			name:  "Ultra compressed flag",
			input: "/user:build --uc",
			checkFlags: func(f *Flags) bool {
				return f.UltraCompressed
			},
		},
		{
			name:  "Evidence flag",
			input: "/user:analyze --evidence",
			checkFlags: func(f *Flags) bool {
				return f.Evidence
			},
		},
		{
			name:  "Think hard flag",
			input: "/user:design --think-hard",
			checkFlags: func(f *Flags) bool {
				return f.Think == "deep"
			},
		},
		{
			name:  "Ultrathink flag",
			input: "/user:troubleshoot --ultrathink",
			checkFlags: func(f *Flags) bool {
				return f.Think == "ultra"
			},
		},
		{
			name:  "Sequential flag",
			input: "/user:load --seq",
			checkFlags: func(f *Flags) bool {
				return f.Sequential
			},
		},
		{
			name:  "All MCP flag",
			input: "/user:spawn --all-mcp",
			checkFlags: func(f *Flags) bool {
				return f.AllMCP
			},
		},
		{
			name:  "Custom flag with value",
			input: "/user:improve --threshold 95%",
			checkFlags: func(f *Flags) bool {
				return f.Additional["threshold"] == "95%"
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsed, err := ParseSuperClaudeCommand(tt.input)
			if err != nil {
				t.Fatalf("Failed to parse command: %v", err)
			}

			if !tt.checkFlags(parsed.Flags) {
				t.Errorf("Flag check failed for input: %s", tt.input)
			}
		})
	}
}

func TestPersonaSelection(t *testing.T) {
	tests := []struct {
		command     string
		wantPersona string
	}{
		{"build", "architect"},
		{"test", "qa"},
		{"scan", "security"},
		{"improve", "refactorer"},
		{"troubleshoot", "analyzer"},
		{"design", "architect"},
		{"deploy", "backend"},
		{"document", "mentor"},
		{"review", "analyzer"},
		{"explain", "mentor"},
	}

	for _, tt := range tests {
		t.Run(tt.command, func(t *testing.T) {
			persona := GetPersonaForCommand(tt.command)
			if persona != tt.wantPersona {
				t.Errorf("GetPersonaForCommand(%s) = %v, want %v", tt.command, persona, tt.wantPersona)
			}
		})
	}
}

func TestGetThinkingTokens(t *testing.T) {
	tests := []struct {
		mode       string
		wantTokens int
	}{
		{"", 4000},
		{"standard", 8000},
		{"deep", 16000},
		{"ultra", 32000},
	}

	for _, tt := range tests {
		t.Run(tt.mode, func(t *testing.T) {
			tokens := GetThinkingTokens(tt.mode)
			if tokens != tt.wantTokens {
				t.Errorf("GetThinkingTokens(%s) = %v, want %v", tt.mode, tokens, tt.wantTokens)
			}
		})
	}
}
