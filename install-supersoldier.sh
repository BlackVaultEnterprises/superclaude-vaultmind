#!/bin/bash
set -e

echo "⚡ Installing SuperSoldierAI (SuperClaude + OpenCode)"
echo "=================================================="
echo ""

# Check prerequisites
echo "🔍 Checking prerequisites..."

# Check Go
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.21+"
    exit 1
fi

# Check Ollama
if ! command -v ollama &> /dev/null; then
    echo "⚠️  Ollama not found. Installing..."
    curl -fsSL https://ollama.com/install.sh | sh
fi

# Apply SuperClaude integration
echo "🔧 Integrating SuperClaude commands..."

# Check if integration is already applied
if ! grep -q "superclaude" internal/tui/page/chat.go 2>/dev/null; then
    echo "📝 Applying SuperClaude integration patch..."
    
    # Create backup
    cp internal/tui/page/chat.go internal/tui/page/chat.go.backup
    
    # Apply integration manually (since patch might fail)
    # Find the sendMessage function and add SuperClaude handler
    cat > /tmp/integrate-superclaude.go << 'EOF'
package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"
)

func main() {
    input, _ := os.Open("internal/tui/page/chat.go")
    defer input.Close()
    
    output, _ := os.Create("internal/tui/page/chat.go.new")
    defer output.Close()
    
    scanner := bufio.NewScanner(input)
    writer := bufio.NewWriter(output)
    
    inSendMessage := false
    importAdded := false
    
    for scanner.Scan() {
        line := scanner.Text()
        
        // Add import if not present
        if !importAdded && strings.Contains(line, "import (") {
            writer.WriteString(line + "\n")
            for scanner.Scan() {
                line = scanner.Text()
                writer.WriteString(line + "\n")
                if strings.Contains(line, `"github.com/opencode-ai/opencode/internal/message"`) {
                    writer.WriteString(`	"github.com/opencode-ai/opencode/internal/superclaude"` + "\n")
                    importAdded = true
                }
                if line == ")" {
                    break
                }
            }
            continue
        }
        
        // Add SuperClaude handler in sendMessage
        if strings.Contains(line, "func (p *chatPage) sendMessage(text string, attachments []message.Attachment) tea.Cmd {") {
            inSendMessage = true
            writer.WriteString(line + "\n")
            writer.WriteString(`	// Check if it's a SuperClaude command
	superHandler := superclaude.NewSuperClaudeHandler(p.app.CoderAgent)
	handled, err := superHandler.HandleCommand(context.Background(), p.session.ID, text)
	if err != nil {
		return util.ReportError(err)
	}
	
	// If it was handled as a SuperClaude command, return early
	if handled {
		// Just ensure we have a session if needed
		if p.session.ID == "" {
			return p.createSessionAndSend(text, attachments)
		}
		return nil
	}
	
`)
            continue
        }
        
        writer.WriteString(line + "\n")
    }
    
    writer.Flush()
    os.Rename("internal/tui/page/chat.go.new", "internal/tui/page/chat.go")
}
EOF
    
    go run /tmp/integrate-superclaude.go
    rm /tmp/integrate-superclaude.go
else
    echo "✅ SuperClaude integration already applied"
fi

# Build the enhanced version
echo "🔨 Building SuperSoldierAI..."
go build -o supersoldier-local

if [ $? -ne 0 ]; then
    echo "❌ Build failed"
    exit 1
fi

# Install to user bin
echo "📦 Installing to ~/.local/bin..."
mkdir -p ~/.local/bin
cp supersoldier-local ~/.local/bin/

# Check if ~/.local/bin is in PATH
if [[ ":$PATH:" != *":$HOME/.local/bin:"* ]]; then
    echo "⚠️  Adding ~/.local/bin to PATH..."
    echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.bashrc
    echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.zshrc 2>/dev/null || true
    echo ""
    echo "📌 Please run: source ~/.bashrc"
fi

# Pull default model if not present
if ! ollama list 2>/dev/null | grep -q "deepseek-coder:6.7b"; then
    echo "📥 Pulling recommended AI model..."
    ollama pull deepseek-coder:6.7b
fi

# Create launch script
cat > ~/.local/bin/supersoldier << 'EOF'
#!/bin/bash
# SuperSoldierAI Launcher

# Ensure Ollama is running
if ! curl -s http://localhost:11434/api/generate -d '{"model":"deepseek-coder:6.7b","prompt":"test"}' > /dev/null 2>&1; then
    echo "Starting Ollama server..."
    nohup ollama serve > /tmp/ollama.log 2>&1 &
    sleep 3
fi

# Launch SuperSoldierAI
exec supersoldier-local "$@"
EOF

chmod +x ~/.local/bin/supersoldier

echo ""
echo "✅ SuperSoldierAI Installation Complete!"
echo ""
echo "🚀 To start using SuperSoldierAI:"
echo "   1. source ~/.bashrc"
echo "   2. supersoldier"
echo ""
echo "📚 Try these SuperClaude commands:"
echo "   /user:analyze --architecture"
echo "   /persona:frontend → /user:build --react"
echo "   /user:spawn --task 'Full Stack Build' --all-mcp"
echo ""
echo "⚡ For maximum deployment: ./SUPERSOLDIER-MAXIMUM-DEPLOYMENT.sh"
echo ""