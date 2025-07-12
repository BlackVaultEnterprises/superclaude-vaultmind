#!/bin/bash

# SuperSoldierAI Build Script
# Builds OpenCode with SuperSoldierAI military-grade enhancements

set -e

echo "⚡ Building SuperSoldierAI System"
echo "================================="
echo ""

# Check Go version
GO_VERSION=$(go version 2>/dev/null | grep -oE '[0-9]+\.[0-9]+' | head -1)
if [ -z "$GO_VERSION" ]; then
    echo "❌ Go is not installed. Please install Go 1.21+"
    exit 1
fi

echo "✅ Go version: $GO_VERSION"

# Apply SuperClaude integration patch if not already applied
if ! grep -q "superclaude" internal/tui/page/chat.go 2>/dev/null; then
    echo "📝 Applying SuperClaude integration..."
    # Note: In production, you'd apply the actual patch
    echo "   (Manual integration required - see superclaude-integration.patch)"
fi

# Build the project
echo "🔨 Building OpenCode with SuperClaude..."
go build -o superclaude-local main.go

if [ $? -eq 0 ]; then
    echo "✅ Build successful!"
    echo ""
    echo "📦 Binary created: ./superclaude-local"
    echo ""
    echo "🚀 Next steps:"
    echo "   1. Ensure Ollama is running: ollama serve"
    echo "   2. Pull a model: ollama pull deepseek-coder:6.7b"
    echo "   3. Run SuperClaude: ./superclaude-local"
    echo ""
    echo "📚 Try these commands:"
    echo "   /user:analyze --architecture"
    echo "   /persona:frontend → /user:build --react"
    echo "   /user:test --coverage --think"
    echo ""
else
    echo "❌ Build failed"
    exit 1
fi