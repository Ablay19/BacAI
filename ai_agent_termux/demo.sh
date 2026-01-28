#!/bin/bash

echo "=== A.I.D.A - Android Intelligent Document Analyst ==="
echo "Demonstration Script"
echo ""

echo "1. Creating sample documents for testing..."
mkdir -p ~/demo_docs

# Create text document
cat > ~/demo_docs/project_plan.txt << EOF
Project: Mobile App Development
Timeline: 3 months
Team: 5 developers, 2 designers
Budget: $50,000
Requirements:
- User authentication
- Real-time messaging
- Payment integration
- Analytics dashboard
EOF

# Create another text document
cat > ~/demo_docs/meeting_notes.txt << EOF
Team Meeting - June 15, 2026
Attendees: Alice, Bob, Charlie, Diana
Topics discussed:
1. Progress on user authentication module
2. Issues with real-time messaging implementation
3. Budget approval for third-party services
Action items:
- Alice: Complete authentication by June 20
- Bob: Resolve messaging latency issues
- Charlie: Prepare budget breakdown
EOF

echo "Created sample documents in ~/demo_docs/"

echo ""
echo "2. Running AI Agent scan on demo documents..."
cd /home/abdoullahelvogani/organizer/ai_agent_termux
go run main.go scan

echo ""
echo "3. Checking processed output..."
ls -la ~/processed_data/

echo ""
echo "4. Showing processed content of one file..."
for file in ~/processed_data/*.json; do
  if [ -f "$file" ]; then
    echo "Processed file: $(basename "$file")"
    cat "$file" | jq '.'
    break
  fi
done

echo ""
echo "=== Demo Complete ==="
echo "The AI Agent has successfully processed the documents."
echo "Check the ~/processed_data/ directory for full results."