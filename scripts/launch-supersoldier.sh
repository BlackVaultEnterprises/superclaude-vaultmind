#!/bin/bash

# SuperSoldierAI Launch Script
# Shows banner and provides quick start instructions

clear

# Display the banner
cat SUPERSOLDIER-BANNER.txt

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Check if Ollama is running
if curl -s http://localhost:11434/api/generate -d '{"model":"deepseek-coder:6.7b","prompt":"test"}' > /dev/null 2>&1; then
    echo "✅ Ollama Status: ONLINE"
else
    echo "❌ Ollama Status: OFFLINE"
    echo "   Run: ollama serve"
fi

# Check if model is available
if ollama list 2>/dev/null | grep -q "deepseek-coder:6.7b"; then
    echo "✅ AI Model: READY (deepseek-coder:6.7b)"
else
    echo "❌ AI Model: NOT FOUND"
    echo "   Run: ollama pull deepseek-coder:6.7b"
fi

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "🎯 QUICK COMMANDS:"
echo ""
echo "  Single Agent Operations:"
echo "    /user:analyze --architecture        # Analyze system architecture"
echo "    /user:build --react --tdd          # Build with TDD"
echo "    /user:test --coverage --e2e        # Run comprehensive tests"
echo "    /user:improve --perf --iterate     # Optimize performance"
echo ""
echo "  Persona Operations:"
echo "    /persona:security → /user:scan --owasp"
echo "    /persona:frontend → /user:build --react --magic"
echo "    /persona:performance → /user:analyze --bottlenecks"
echo ""
echo "  Maximum Deployment:"
echo "    ./SUPERSOLDIER-MAXIMUM-DEPLOYMENT.sh    # Deploy all 10 battalions"
echo ""
echo "  Mission Control:"
echo "    cat SUPERSOLDIER-MISSION-CONTROL.md     # View battalion status"
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "⚡ SUPERSOLDIERAI READY FOR DEPLOYMENT ⚡"
echo ""