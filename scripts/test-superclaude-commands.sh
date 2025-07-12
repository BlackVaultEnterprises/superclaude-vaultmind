#!/bin/bash

# SuperClaude Command Test Suite
# Tests all 18 commands with various flags and scenarios

set -e

echo "🧪 SuperClaude Command Test Suite"
echo "================================="

# Check if superclaude binary exists
if [ ! -f "./superclaude" ]; then
    echo "❌ Error: superclaude binary not found!"
    echo "Building..."
    go build -o superclaude || exit 1
fi

# Test function
test_command() {
    local cmd="$1"
    local description="$2"
    echo -n "Testing: $description... "
    
    # Create test input file
    echo "$cmd" > test_input.txt
    
    # Run command with timeout
    if timeout 5s ./superclaude --provider openrouter --model mistralai/mixtral-8x7b-instruct --non-interactive < test_input.txt > /dev/null 2>&1; then
        echo "✅ PASS"
    else
        echo "❌ FAIL"
        FAILED=$((FAILED + 1))
    fi
    
    rm -f test_input.txt
}

# Initialize counters
FAILED=0

echo ""
echo "1️⃣ Testing Basic Commands"
echo "--------------------------"
test_command "/user:analyze ." "Analyze current directory"
test_command "/user:build --react" "Build React app"
test_command "/user:test" "Generate tests"
test_command "/user:improve" "Improve code"
test_command "/user:design api" "Design API"
test_command "/user:review" "Code review"

echo ""
echo "2️⃣ Testing Commands with Flags"
echo "--------------------------------"
test_command "/user:analyze . --deep" "Deep analysis"
test_command "/user:build --react --typescript" "Build with TypeScript"
test_command "/user:test --coverage --e2e" "Test with coverage"
test_command "/user:improve --performance" "Performance optimization"
test_command "/user:scan --owasp" "Security scan"

echo ""
echo "3️⃣ Testing Thinking Modes"
echo "---------------------------"
test_command "/user:design system --think-hard" "Design with deep thinking"
test_command "/user:analyze . --ultrathink" "Ultra thinking analysis"
test_command "/user:build app --uc" "Ultra compressed build"

echo ""
echo "4️⃣ Testing Persona Overrides"
echo "------------------------------"
test_command "/persona:security → /user:review" "Security review"
test_command "/persona:frontend → /user:build" "Frontend build"
test_command "/persona:architect → /user:design" "Architect design"

echo ""
echo "5️⃣ Testing Complex Scenarios"
echo "------------------------------"
test_command "/user:spawn --parallel analyze,test,review" "Parallel spawn"
test_command "/user:git commit --message='test'" "Git operations"
test_command "/user:migrate database --plan" "Migration planning"

echo ""
echo "6️⃣ Testing All Commands"
echo "------------------------"
commands=(
    "analyze" "build" "test" "improve" "design" "review"
    "deploy" "document" "troubleshoot" "migrate" "cleanup"
    "explain" "estimate" "scan" "git" "spawn" "dev-setup" "load"
)

for cmd in "${commands[@]}"; do
    test_command "/user:$cmd" "Command: $cmd"
done

echo ""
echo "================================="
echo "Test Results:"
echo "Total Commands: 18"
echo "Failed: $FAILED"

if [ $FAILED -eq 0 ]; then
    echo "✅ All tests passed!"
    exit 0
else
    echo "❌ $FAILED tests failed!"
    exit 1
fi