package superclaude

import (
	"fmt"
	"strings"
	"text/template"
)

// SuperClaudeCommand represents a SuperClaude command with its configuration
type SuperClaudeCommand struct {
	Name        string
	Persona     string
	Flags       map[string]string
	Template    string
	Description string
}

// CommandTemplate holds the template structure for commands
type CommandTemplate struct {
	Command         string
	Target          string
	Persona         string
	Flags           map[string]interface{}
	UltraCompressed bool
	Think           bool
	ThinkLevel      string
	AnalysisType    string
	Evidence        bool
}

// Commands defines all available SuperClaude commands
var Commands = map[string]SuperClaudeCommand{
	"build": {
		Name:        "build",
		Description: "Build a component or system with specified technology",
		Template: `You are a {{.Persona}} specialist. Build {{.Target}} using {{.Flags}}.
        
CRITICAL RULES:
- Evidence-based decisions only (cite sources)
- {{if .UltraCompressed}}Use 70% fewer tokens{{end}}
- {{if .Think}}Provide {{.ThinkLevel}} analysis depth{{end}}
- Follow TDD practices if --tdd flag is set
- Include comprehensive error handling
        
Execute: {{.Command}}`,
	},

	"analyze": {
		Name:        "analyze",
		Description: "Analyze code, architecture, or system components",
		Template: `Analyze {{.Target}} as {{.Persona}}.
        
Focus: {{.AnalysisType}}
Depth: {{.ThinkLevel}}
Output: Evidence-based findings with citations
{{if .Evidence}}Include: External documentation references{{end}}`,
	},

	"test": {
		Name:        "test",
		Description: "Run tests with specified coverage and frameworks",
		Template: `Execute {{.Target}} tests as {{.Persona}}.
		
Test Type: {{.Flags.type}}
Coverage Target: {{.Flags.coverage}}
{{if .Flags.watch}}Watch Mode: Enabled{{end}}
Report Format: Detailed with metrics`,
	},

	"improve": {
		Name:        "improve",
		Description: "Improve code quality, performance, or architecture",
		Template: `Improve {{.Target}} as {{.Persona}}.

IMPROVEMENT FOCUS:
{{if .Flags.quality}}Code Quality: Naming clarity | Function extraction | Duplication removal | Complexity reduction{{end}}
{{if .Flags.perf}}Performance: Algorithm optimization | Query optimization | Caching strategies | Memory efficiency{{end}}
{{if .Flags.arch}}Architecture: Design patterns | Dependency injection | Layer separation | Scalability patterns{{end}}
{{if .Flags.refactor}}Refactoring: Safe changes preserving behavior{{end}}

PROCESS:
1. Analysis: Current state assessment | Identify improvement areas
2. Planning: Safe refactoring path | Preserve functionality
3. Implementation: Small atomic changes | Continuous testing
4. Validation: Behavior preservation | Performance gains

{{if .Flags.threshold}}Quality Threshold: {{.Flags.threshold}}{{end}}
{{if .Flags.iterate}}Iterative Mode: Continue until threshold met{{end}}
{{if .Flags.metrics}}Metrics: Show before/after measurements{{end}}
{{if .Flags.safe}}Safe Mode: Conservative changes only{{end}}

Validation: Evidence-based improvements with measurable impact`,
	},

	"troubleshoot": {
		Name:        "troubleshoot",
		Description: "Debug and resolve issues systematically",
		Template: `Troubleshoot {{.Target}} as {{.Persona}}.
		
Method: {{.Flags.method}}
{{if .Flags.investigate}}Deep Investigation: Enabled{{end}}
{{if .Flags.fiveWhys}}Five Whys Analysis: Enabled{{end}}`,
	},

	"design": {
		Name:        "design",
		Description: "Design systems, APIs, or architectures",
		Template: `Design {{.Target}} as {{.Persona}}.
		
Patterns: {{.Flags.patterns}}
Architecture: {{.Flags.architecture}}
{{if .Evidence}}Documentation: Include references{{end}}`,
	},

	"deploy": {
		Name:        "deploy",
		Description: "Deploy applications or services",
		Template: `Deploy {{.Target}} to {{.Flags.env}} as {{.Persona}}.
		
Environment: {{.Flags.env}}
{{if .Flags.dryRun}}Dry Run Mode: Enabled{{end}}
{{if .Flags.rollback}}Rollback Ready: Yes{{end}}`,
	},

	"scan": {
		Name:        "scan",
		Description: "Scan for security, quality, or compliance issues",
		Template: `Scan {{.Target}} for {{.Flags.type}} as {{.Persona}}.
		
Scan Type: {{.Flags.type}}
{{if .Flags.owasp}}OWASP Standards: Applied{{end}}
{{if .Flags.validate}}Validation: Strict{{end}}`,
	},

	"document": {
		Name:        "document",
		Description: "Generate documentation for code or systems",
		Template: `Document {{.Target}} as {{.Persona}}.
		
Type: {{.Flags.type}}
Format: {{.Flags.format}}
{{if .Flags.interactive}}Interactive Mode: Enabled{{end}}`,
	},

	"review": {
		Name:        "review",
		Description: "Review code, architecture, or documentation",
		Template: `Review {{.Target}} as {{.Persona}}.
		
Focus: {{.Flags.focus}}
{{if .Evidence}}Evidence Required: Yes{{end}}
Standards: Comprehensive analysis`,
	},

	"migrate": {
		Name:        "migrate",
		Description: "Migrate data, code, or systems",
		Template: `Migrate {{.Target}} as {{.Persona}}.
		
{{if .Flags.dryRun}}Dry Run: Enabled{{end}}
{{if .Flags.validate}}Validation: Enabled{{end}}
{{if .Flags.rollback}}Rollback Strategy: Prepared{{end}}`,
	},

	"cleanup": {
		Name:        "cleanup",
		Description: "Clean up code, dependencies, or resources",
		Template: `Cleanup {{.Target}} as {{.Persona}}.
		
Scope: {{.Flags.scope}}
{{if .Flags.optimize}}Optimization: Enabled{{end}}
{{if .Flags.validate}}Validation: Post-cleanup{{end}}`,
	},

	"explain": {
		Name:        "explain",
		Description: "Explain code, concepts, or systems",
		Template: `Explain {{.Target}} as {{.Persona}}.
		
Depth: {{.Flags.depth}}
{{if .Flags.visual}}Visual Aids: Include diagrams{{end}}
{{if .Flags.examples}}Examples: Provide practical examples{{end}}`,
	},

	"estimate": {
		Name:        "estimate",
		Description: "Estimate time, resources, or complexity",
		Template: `Estimate {{.Target}} as {{.Persona}}.
		
Type: {{.Flags.type}}
Confidence: Evidence-based
Include: Risk factors and assumptions`,
	},

	"dev-setup": {
		Name:        "dev-setup",
		Description: "Set up development environment",
		Template: `Setup development environment for {{.Target}} as {{.Persona}}.
		
Components: {{.Flags.components}}
{{if .Flags.ci}}CI/CD: Configure{{end}}
{{if .Flags.automation}}Automation: Enable{{end}}`,
	},

	"load": {
		Name:        "load",
		Description: "Load and analyze codebase or data",
		Template: `Load {{.Target}} as {{.Persona}}.
		
Depth: {{.Flags.depth}}
{{if .Flags.seq}}Sequential Processing: Enabled{{end}}
Analysis: Comprehensive understanding`,
	},

	"git": {
		Name:        "git",
		Description: "Execute git operations",
		Template: `Execute git {{.Target}} as {{.Persona}}.
		
Operation: {{.Target}}
{{if .Flags.checkpoint}}Checkpoint: Create{{end}}
Safety: Validate before execution`,
	},

	"spawn": {
		Name:        "spawn",
		Description: "Spawn multi-agent tasks",
		Template: `Spawn multi-agent task {{.Target}} as {{.Persona}}.
		
Agents: {{.Flags.agents}}
{{if .Flags.parallel}}Parallel Execution: Enabled{{end}}
Coordination: Managed workflow`,
	},
}

// BuildPrompt generates the final prompt from a command and context
func (cmd *SuperClaudeCommand) BuildPrompt(persona Persona, flags *Flags, target string, rawCommand string) (string, error) {
	tmpl, err := template.New("command").Parse(cmd.Template)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	data := CommandTemplate{
		Command:         rawCommand,
		Target:          target,
		Persona:         persona.Name,
		Flags:           convertFlagsToMap(flags),
		UltraCompressed: flags.UltraCompressed,
		Think:           flags.Think != "",
		ThinkLevel:      flags.Think,
		Evidence:        flags.Evidence,
	}

	if analysisType, ok := flags.Additional["type"]; ok {
		data.AnalysisType = analysisType
	}

	var result strings.Builder

	// Add persona context
	result.WriteString(fmt.Sprintf("IDENTITY: %s\n", persona.Identity))
	result.WriteString(fmt.Sprintf("CORE BELIEF: %s\n", persona.CoreBelief))
	result.WriteString(fmt.Sprintf("DECISION FRAMEWORK: %s\n", persona.DecisionFramework))
	result.WriteString(fmt.Sprintf("COMMUNICATION STYLE: %s\n\n", persona.CommunicationStyle))

	// Execute template
	if err := tmpl.Execute(&result, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return result.String(), nil
}

// convertFlagsToMap converts Flags struct to map for template use
func convertFlagsToMap(flags *Flags) map[string]interface{} {
	result := make(map[string]interface{})

	// Add all additional flags
	for k, v := range flags.Additional {
		result[k] = v
	}

	// Add specific flags
	if flags.Think != "" {
		result["think"] = flags.Think
	}
	if flags.UltraCompressed {
		result["uc"] = true
	}
	if flags.Plan {
		result["plan"] = true
	}
	if flags.Evidence {
		result["evidence"] = true
	}
	if flags.ValidationOnly {
		result["validate"] = true
	}

	return result
}
