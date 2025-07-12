#!/bin/bash

# SuperClaude Local Demo Script
# This demonstrates various SuperClaude commands and features

echo "🚀 SuperClaude Local Demo"
echo "========================="
echo ""

# Ensure Ollama is running
if ! curl -s http://localhost:11434/api/generate -d '{"model":"deepseek-coder:6.7b","prompt":"test"}' > /dev/null 2>&1; then
    echo "❌ Ollama is not running. Starting Ollama..."
    ollama serve &
    sleep 5
fi

echo "✅ Ollama is running"
echo ""

# Function to simulate command execution
run_command() {
    local cmd="$1"
    local desc="$2"
    
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "📌 $desc"
    echo "💻 Command: $cmd"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo ""
    
    # In real usage, this would be sent to the OpenCode/SuperClaude interface
    echo "✨ This command would execute with the following behavior:"
    
    case "$cmd" in
        *"analyze"*)
            echo "  - Analyzing codebase structure"
            echo "  - Identifying patterns and anti-patterns"
            echo "  - Generating architectural insights"
            ;;
        *"build"*)
            echo "  - Creating component structure"
            echo "  - Implementing with specified framework"
            echo "  - Following TDD practices if specified"
            ;;
        *"test"*)
            echo "  - Running test suite"
            echo "  - Generating coverage reports"
            echo "  - Identifying missing test cases"
            ;;
        *"improve"*)
            echo "  - Analyzing performance bottlenecks"
            echo "  - Implementing optimizations"
            echo "  - Validating improvements"
            ;;
        *"scan"*)
            echo "  - Scanning for security vulnerabilities"
            echo "  - Checking OWASP compliance"
            echo "  - Generating security report"
            ;;
    esac
    
    echo ""
    sleep 2
}

# Demo 1: Basic Commands
echo "📘 DEMO 1: Basic SuperClaude Commands"
echo "======================================"
echo ""

run_command "/user:analyze --architecture" \
    "Analyze project architecture"

run_command "/user:build --react --component TodoList" \
    "Build a React component"

run_command "/user:test --coverage --unit" \
    "Run unit tests with coverage"

# Demo 2: Persona Usage
echo "📘 DEMO 2: Using Cognitive Personas"
echo "===================================="
echo ""

run_command "/persona:security → /user:scan --owasp --deps" \
    "Security scan with security persona"

run_command "/persona:frontend → /user:build --react --magic --uc" \
    "Frontend build with UI specialist"

run_command "/persona:performance → /user:improve --cache --threshold 95%" \
    "Performance optimization"

# Demo 3: Advanced Flags
echo "📘 DEMO 3: Advanced Flags and Thinking Modes"
echo "============================================"
echo ""

run_command "/user:design api --ddd --think-hard --evidence" \
    "Design API with deep thinking and evidence"

run_command "/user:troubleshoot --investigate --five-whys --ultrathink" \
    "Ultra-deep troubleshooting analysis"

run_command "/user:analyze --code --uc --seq" \
    "Ultra-compressed sequential analysis"

# Demo 4: Multi-Agent Orchestration
echo "📘 DEMO 4: Multi-Agent Orchestration"
echo "===================================="
echo ""

run_command '/user:spawn --task "Full Stack Build" --all-mcp' \
    "Spawn multiple agents for full-stack development"

# Demo 5: Complex Workflow
echo "📘 DEMO 5: Complex Development Workflow"
echo "======================================="
echo ""

echo "🔄 Simulating a complete feature development cycle:"
echo ""

run_command "/persona:architect → /user:design --api --microservices --plan" \
    "Step 1: Architecture design"

run_command "/persona:backend → /user:build --api --tdd --openapi" \
    "Step 2: Backend implementation"

run_command "/persona:frontend → /user:build --react --tdd --watch" \
    "Step 3: Frontend development"

run_command "/persona:qa → /user:test --e2e --coverage --pup" \
    "Step 4: Comprehensive testing"

run_command "/persona:security → /user:scan --security --validate" \
    "Step 5: Security validation"

run_command "/user:deploy --env staging --dry-run --plan" \
    "Step 6: Deployment preparation"

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✅ Demo Complete!"
echo ""
echo "📚 Key Takeaways:"
echo "  • SuperClaude commands start with /user: or /persona:"
echo "  • 9 specialized personas for different tasks"
echo "  • Thinking modes: --think, --think-hard, --ultrathink"
echo "  • Evidence-based operation with --evidence flag"
echo "  • Ultra-compressed mode with --uc for efficiency"
echo "  • Multi-agent orchestration with /user:spawn"
echo ""
echo "🚀 Ready to supercharge your development with SuperClaude!"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"