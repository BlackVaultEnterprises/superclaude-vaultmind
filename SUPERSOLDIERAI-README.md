# SuperSoldierAI - Maximum Agent Intelligence System

An elite AI command framework that transforms OpenCode into a multi-agent orchestration platform with SuperClaude's advanced capabilities, military-grade coordination, and evidence-based operation.

## 🚀 Quick Start

### Prerequisites
- Go 1.21+
- Ollama installed and running
- 8GB+ RAM for local models

### Installation

```bash
# 1. Clone and build
git clone https://github.com/opencode-ai/opencode superclaude-local
cd superclaude-local
go build

# 2. Install Ollama (if not already installed)
curl -fsSL https://ollama.com/install.sh | sh

# 3. Pull recommended models
ollama pull deepseek-coder:6.7b  # Balanced performance
ollama pull codellama:13b        # Alternative option

# 4. Start Ollama server
ollama serve
```

## 🎯 SuperClaude Commands

All SuperClaude commands follow this format:
```
/user:command [target] [--flags]
```

Or with persona override:
```
/persona:architect → /user:command [target] [--flags]
```

### Available Commands

| Command | Description | Example |
|---------|-------------|---------|
| `build` | Build components with specified technology | `/user:build --react --tdd` |
| `analyze` | Analyze code, architecture, or systems | `/user:analyze --architecture --deep` |
| `test` | Run tests with coverage | `/user:test --coverage --e2e` |
| `improve` | Improve code quality or performance | `/user:improve --performance --threshold 95%` |
| `troubleshoot` | Debug issues systematically | `/user:troubleshoot --investigate --five-whys` |
| `design` | Design systems or APIs | `/user:design --api --ddd --openapi` |
| `deploy` | Deploy applications | `/user:deploy --env staging --dry-run` |
| `scan` | Security and quality scanning | `/user:scan --security --owasp` |
| `document` | Generate documentation | `/user:document --api --interactive` |
| `review` | Review code or architecture | `/user:review --quality --evidence` |
| `migrate` | Migrate data or systems | `/user:migrate --validate --rollback` |
| `cleanup` | Clean up code or dependencies | `/user:cleanup --all --optimize` |
| `explain` | Explain code or concepts | `/user:explain --depth expert --visual` |
| `estimate` | Estimate time or complexity | `/user:estimate --detailed` |
| `dev-setup` | Set up development environment | `/user:dev-setup --ci --automation` |
| `load` | Load and analyze codebase | `/user:load --depth deep --seq` |
| `git` | Git operations | `/user:git --checkpoint v1.0` |
| `spawn` | Spawn multi-agent tasks | `/user:spawn --task "full-build" --all-mcp` |

### Universal Flags

| Flag | Description |
|------|-------------|
| `--persona-[name]` | Override default persona |
| `--think` | Standard thinking mode |
| `--think-hard` | Deep analysis mode |
| `--ultrathink` | Exhaustive analysis mode |
| `--uc` | Ultra-compressed output (70% fewer tokens) |
| `--evidence` / `--c7` | Require evidence for all claims |
| `--seq` | Sequential processing |
| `--plan` | Create execution plan |
| `--validate` | Validation only mode |
| `--all-mcp` | Use all available MCP tools |

## 🧠 Cognitive Personas

### Available Personas

1. **architect** - Systems design, scalability, patterns
2. **frontend** - UI/UX, React, performance
3. **backend** - APIs, databases, distributed systems
4. **security** - Security analysis, threat modeling
5. **qa** - Testing, quality assurance, automation
6. **refactorer** - Code cleanup, technical debt
7. **performance** - Optimization, benchmarking
8. **analyzer** - Code analysis, complexity detection
9. **mentor** - Documentation, teaching, examples

### Using Personas

```bash
# Explicit persona selection
/persona:security → /user:scan --owasp

# Flag-based selection
/user:build --react --persona-frontend

# Default persona (auto-selected based on command)
/user:test  # Uses 'qa' persona
```

## 🔧 Configuration

Edit `~/.config/superclaude/config.yaml`:

```yaml
providers:
  primary: "ollama"
  
  ollama:
    models:
      fast: "codellama:7b"
      balanced: "deepseek-coder:6.7b"
      powerful: "codellama:34b"

personas:
  default: "architect"
  
commands:
  evidence_required: true
  auto_compress: true
```

## 💡 Example Workflows

### Full-Stack Development
```bash
# Architecture design
/persona:architect → /user:design --api --ddd --microservices

# Backend implementation
/persona:backend → /user:build --api --tdd --openapi

# Frontend development
/persona:frontend → /user:build --react --magic --tdd

# Quality assurance
/persona:qa → /user:test --coverage --e2e --integration
```

### Performance Optimization
```bash
# Analyze bottlenecks
/persona:analyzer → /user:analyze --performance --deep

# Profile and investigate
/persona:performance → /user:troubleshoot --perf --five-whys

# Implement improvements
/user:improve --performance --cache --threshold 95%
```

### Security Review
```bash
# Security scan
/persona:security → /user:scan --security --owasp --deps

# Threat analysis
/user:troubleshoot --investigate --security --evidence

# Implement fixes
/user:improve --security --compliance --strict
```

## 🚀 Advanced Features

### Multi-Agent Orchestration
```bash
/user:spawn --task "production-deployment" --all-mcp --ultrathink
```

This spawns multiple specialized agents working in parallel:
- Architecture validation
- Security scanning
- Performance testing
- Deployment preparation

### Evidence-Based Operation
```bash
/user:build --react --evidence
```

Forces all decisions to be backed by:
- Documentation references
- Best practice citations
- Performance benchmarks
- Security standards

### Thinking Modes
```bash
# Quick analysis
/user:analyze --code

# Thorough analysis
/user:analyze --code --think

# Deep analysis
/user:analyze --code --think-hard

# Exhaustive analysis
/user:analyze --code --ultrathink
```

## 🛠️ Development

### Building from Source
```bash
go build -o superclaude-local
./superclaude-local
```

### Adding New Commands
1. Add command definition to `internal/superclaude/commands.go`
2. Update command routing in `integration.go`
3. Add tests in `superclaude_test.go`

### Adding New Personas
1. Define persona in `internal/superclaude/personas.go`
2. Map commands to persona in `GetPersonaForCommand()`
3. Update documentation

## 📊 Performance Tips

1. **Model Selection**
   - Use `fast` model for simple tasks
   - Use `balanced` for most development
   - Reserve `powerful` for complex architecture

2. **Token Optimization**
   - Use `--uc` for routine tasks
   - Combine with `--evidence` for concise, factual output
   - Use thinking modes sparingly for complex problems

3. **Caching**
   - Results cached for 15 minutes by default
   - Adjust `cache_ttl` in config for different needs

## 🐛 Troubleshooting

### Ollama Connection Issues
```bash
# Check if Ollama is running
curl http://localhost:11434/api/generate -d '{
  "model": "deepseek-coder:6.7b",
  "prompt": "Hello"
}'

# Restart Ollama
pkill ollama
ollama serve
```

### Model Performance
```bash
# List available models
ollama list

# Remove unused models
ollama rm model-name

# Check model info
ollama show deepseek-coder:6.7b
```

## 🤝 Contributing

1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing`)
5. Open Pull Request

## 📄 License

This project extends OpenCode and maintains its original license while adding the SuperClaude enhancement layer.

---

**SuperClaude Local** - Bringing Claude Code's power to your local development environment 🚀