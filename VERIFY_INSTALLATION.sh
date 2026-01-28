#!/bin/bash

echo "=== A.I.D.A. Installation Verification ==="
echo ""

echo "1. Checking project structure..."
if [ -d "/home/abdoullahelvogani/organizer/ai_agent_termux" ]; then
    echo "✅ Project directory exists"
else
    echo "❌ Project directory missing"
    exit 1
fi

echo ""
echo "2. Checking Go binary..."
if [ -f "/home/abdoullahelvogani/organizer/ai_agent_termux/ai_agent" ]; then
    echo "✅ AI Agent binary exists"
else
    echo "❌ AI Agent binary missing"
    exit 1
fi

echo ""
echo "3. Checking Go modules..."
cd /home/abdoullahelvogani/organizer/ai_agent_termux
if go mod verify >/dev/null 2>&1; then
    echo "✅ Go modules verified"
else
    echo "❌ Go module verification failed"
    exit 1
fi

echo ""
echo "4. Testing AI Agent functionality..."
if ./ai_agent --help 2>/dev/null | grep -q "Usage"; then
    echo "✅ AI Agent help command works"
else
    echo "❌ AI Agent help command failed"
    exit 1
fi

echo ""
echo "5. Checking processed output..."
if [ -d "/home/abdoullahelvogani/processed_data" ] && [ "$(ls -A /home/abdoullahelvogani/processed_data)" ]; then
    echo "✅ Processed output exists"
else
    echo "❌ Processed output missing"
    exit 1
fi

echo ""
echo "========================================="
echo "✅ A.I.D.A. Installation Verified Successfully!"
echo "🚀 Ready for document analysis on Android!"
echo "========================================="