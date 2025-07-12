#!/bin/bash

# Set OpenRouter API key from environment
# export OPENROUTER_API_KEY="your-api-key-here"
if [ -z "$OPENROUTER_API_KEY" ]; then
    echo "❌ Error: OPENROUTER_API_KEY environment variable not set!"
    echo "Please set your API key: export OPENROUTER_API_KEY='your-key'"
    exit 1
fi

# Check if binary exists
if [ ! -f "./superclaude" ]; then
    echo "❌ Error: superclaude binary not found!"
    echo "Please run: go build -o superclaude"
    exit 1
fi

# Launch SuperClaude with OpenRouter
echo "╔══════════════════════════════════════════════════════════════╗"
echo "║            🚀 SuperClaude AI Assistant v1.0                  ║"
echo "╚══════════════════════════════════════════════════════════════╝"
echo ""
echo "Provider: OpenRouter"
echo "Model: ${1:-mistralai/mixtral-8x7b-instruct}"
echo ""
echo "📚 Available Commands:"
echo "  /user:build     - Build projects and features"
echo "  /user:analyze   - Analyze codebases"
echo "  /user:test      - Write and run tests"
echo "  /user:improve   - Optimize and refactor code"
echo "  /user:design    - Design systems and APIs"
echo "  /user:review    - Code review and security audit"
echo ""
echo "🎯 Flags: --uc (compressed) --think-hard (deep analysis)"
echo "══════════════════════════════════════════════════════════════"
echo ""

# Run with the specified model or default
./superclaude --provider openrouter --model "${1:-mistralai/mixtral-8x7b-instruct}"