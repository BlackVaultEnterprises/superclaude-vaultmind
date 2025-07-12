# SuperClaude AI Assistant

[![Go Version](https://img.shields.io/badge/Go-1.24.0-blue.svg)](https://golang.org/dl/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Build Status](https://img.shields.io/github/actions/workflow/status/yourusername/superclaude/ci.yml?branch=main)](https://github.com/yourusername/superclaude/actions)

SuperClaude is an enhanced AI coding assistant that extends OpenCode with 18 specialized commands, 9 cognitive personas, and advanced features for software development.

## 🚀 Features

### 18 Specialized Commands
- **`/user:build`** - Build projects with any framework
- **`/user:analyze`** - Deep codebase analysis
- **`/user:test`** - Comprehensive test generation
- **`/user:improve`** - Code optimization and refactoring
- **`/user:design`** - System architecture and API design
- **`/user:review`** - Code review and security audits
- **`/user:deploy`** - Deployment configuration
- **`/user:document`** - Generate documentation
- **`/user:troubleshoot`** - Debug and fix issues
- **`/user:migrate`** - Database and code migrations
- **`/user:cleanup`** - Code cleanup and tech debt
- **`/user:explain`** - Code explanation
- **`/user:estimate`** - Project estimation
- **`/user:scan`** - Security vulnerability scanning
- **`/user:git`** - Git workflow automation
- **`/user:spawn`** - Multi-agent orchestration
- **`/user:dev-setup`** - Development environment setup
- **`/user:load`** - Load and analyze data

### 9 Cognitive Personas
- **Architect** - System design and scalability
- **Frontend** - UI/UX and React expertise
- **Backend** - APIs and distributed systems
- **Security** - Security analysis and compliance
- **QA** - Testing and quality assurance
- **Performance** - Optimization specialist
- **Refactorer** - Code quality improvement
- **Analyzer** - Code complexity analysis
- **Mentor** - Documentation and teaching

### Advanced Features
- **Multi-provider support** - OpenRouter, OpenAI, Anthropic, Ollama (local), and more
- **Thinking modes** - Standard, deep (--think-hard), and ultra (--ultrathink)
- **Ultra-compressed mode** (--uc) - Optimized token usage
- **Evidence-based** (--evidence) - Requires concrete proof
- **Session management** - Persistent conversation history
- **IDE integration** - VSCode, Vim, Emacs support

## 📦 Installation

### Option 1: Pre-built Binary
```bash
# Download latest release
curl -L https://github.com/yourusername/superclaude/releases/latest/download/superclaude-linux-amd64 -o superclaude
chmod +x superclaude
sudo mv superclaude /usr/local/bin/
```

### Option 2: Build from Source
```bash
# Clone repository
git clone https://github.com/yourusername/superclaude.git
cd superclaude

# Build
go build -o superclaude

# Install
sudo cp superclaude /usr/local/bin/
```

### Option 3: Docker
```bash
docker pull yourusername/superclaude:latest
docker run -it yourusername/superclaude
```

## 🔧 Configuration

### Environment Variables
```bash
# OpenRouter (recommended)
export OPENROUTER_API_KEY="your-key-here"

# Or other providers
export OPENAI_API_KEY="your-key"
export ANTHROPIC_API_KEY="your-key"
```

### Config File (~/.superclaude/config.json)
```json
{
  "provider": "openrouter",
  "model": "mistralai/mixtral-8x7b-instruct",
  "theme": "gruvbox",
  "editor": "vim",
  "workingDirectory": "~/projects"
}
```

## 🎯 Usage

### Basic Commands
```bash
# Launch with default settings
superclaude

# Specify provider and model
superclaude --provider openrouter --model mistralai/mixtral-8x7b-instruct

# Use local model with Ollama
superclaude --provider ollama --model deepseek-coder:6.7b
```

### SuperClaude Commands
```bash
# Build a React app
/user:build --react --typescript a todo app with drag and drop

# Analyze codebase
/user:analyze codebase --deep

# Write tests with coverage
/user:test --coverage --e2e

# Optimize performance
/user:improve --performance --benchmark

# Security scan
/user:scan --owasp --penetration
```

### Command Flags
- `--uc` - Ultra-compressed responses
- `--think-hard` - Deep analysis mode
- `--ultrathink` - Maximum analysis (32k tokens)
- `--evidence` - Require proof for claims
- `--plan` - Show plan before executing
- `--sequential` - Step-by-step execution
- `--all-mcp` - Use all available MCP tools

### Persona Override
```bash
# Force specific persona
/persona:security → /user:review code

# Combine personas
/persona:architect → /user:design api --think-hard
```

## 🔌 IDE Integration

### VSCode Extension
```bash
# Install extension
code --install-extension superclaude.vscode-superclaude

# Configure in settings.json
{
  "superclaude.apiKey": "your-openrouter-key",
  "superclaude.model": "mistralai/mixtral-8x7b-instruct",
  "superclaude.keybinding": "ctrl+shift+s"
}
```

### Vim Plugin
```vim
" Add to .vimrc
Plug 'superclaude/vim-superclaude'

" Configure
let g:superclaude_api_key = $OPENROUTER_API_KEY
let g:superclaude_model = 'mistralai/mixtral-8x7b-instruct'

" Keybindings
nmap <leader>sc :SuperClaude<CR>
nmap <leader>sb :SuperClaudeBuild<CR>
nmap <leader>st :SuperClaudeTest<CR>
```

### Emacs Package
```elisp
;; Add to init.el
(use-package superclaude
  :ensure t
  :config
  (setq superclaude-api-key (getenv "OPENROUTER_API_KEY"))
  (setq superclaude-model "mistralai/mixtral-8x7b-instruct")
  :bind (("C-c s c" . superclaude)
         ("C-c s b" . superclaude-build)
         ("C-c s t" . superclaude-test)))
```

### JetBrains Plugin
Available in JetBrains Marketplace for IntelliJ IDEA, WebStorm, PyCharm, etc.

## 🧪 Testing

### Run Tests
```bash
# All tests
go test ./...

# With coverage
go test -cover ./...

# Specific package
go test ./internal/superclaude/...

# Benchmarks
go test -bench=. ./...
```

### Test SuperClaude Commands
```bash
# Test script included
./test-superclaude-commands.sh

# Manual testing
./superclaude
> /user:test --dry-run
> /user:analyze . --validate-only
```

## 🏗️ Architecture

```
superclaude/
├── cmd/                    # CLI commands
├── internal/
│   ├── superclaude/       # Core SuperClaude logic
│   │   ├── commands.go    # 18 command definitions
│   │   ├── personas.go    # 9 persona definitions
│   │   ├── flags.go       # Command parsing
│   │   └── integration.go # OpenCode integration
│   ├── llm/               # LLM providers
│   ├── tui/               # Terminal UI
│   └── app/               # Application core
├── ide/                   # IDE integrations
│   ├── vscode/           # VSCode extension
│   ├── vim/              # Vim plugin
│   └── emacs/            # Emacs package
└── scripts/              # Build and deployment
```

## 🚀 Deployment

### Local Installation
```bash
# Build and install
make install

# Or manually
go build -o superclaude
sudo cp superclaude /usr/local/bin/
```

### Docker Deployment
```bash
# Build image
docker build -t superclaude .

# Run container
docker run -it -e OPENROUTER_API_KEY=$OPENROUTER_API_KEY superclaude
```

### Kubernetes
```bash
# Apply manifests
kubectl apply -f k8s/

# Or use Helm
helm install superclaude ./helm/superclaude
```

## 🤝 Contributing

1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push branch (`git push origin feature/amazing`)
5. Open Pull Request

### Development Setup
```bash
# Clone repo
git clone https://github.com/yourusername/superclaude
cd superclaude

# Install dependencies
go mod download

# Run tests
make test

# Build
make build
```

## 📄 License

MIT License - see [LICENSE](LICENSE) file

## 🙏 Acknowledgments

- Built on top of [OpenCode](https://github.com/opencode-ai/opencode)
- Inspired by military cognitive enhancement programs
- Community contributors and testers

## 📞 Support

- **Issues**: [GitHub Issues](https://github.com/yourusername/superclaude/issues)
- **Discussions**: [GitHub Discussions](https://github.com/yourusername/superclaude/discussions)
- **Security**: security@superclaude.ai

---

**SuperClaude** - Enhanced AI Coding Assistant | Free Forever with OpenRouter