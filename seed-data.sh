#!/bin/bash

# Seed mock data into the database via GraphQL mutations
# This script reads mockdata.json and creates customers using createCustomer mutation

set -e

GRAPHQL_URL="${GRAPHQL_URL:-http://localhost:8080/query}"
MOCKDATA_FILE="${MOCKDATA_FILE:-mockdata.json}"

echo "🌱 Seeding database with mock data..."
echo "📍 GraphQL endpoint: $GRAPHQL_URL"

# Check if jq is installed
if ! command -v jq &> /dev/null; then
    echo "❌ Error: jq is required but not installed."
    echo "   Install with: brew install jq"
    exit 1
fi

# Check if mockdata file exists
if [ ! -f "$MOCKDATA_FILE" ]; then
    echo "❌ Error: $MOCKDATA_FILE not found"
    exit 1
fi

# Check if server is running
if ! curl -s -X POST "$GRAPHQL_URL" -d '{"query":"{ __typename }"}' > /dev/null 2>&1; then
    echo "❌ Error: GraphQL server not responding at $GRAPHQL_URL"
    echo "   Make sure the server is running: make dev"
    exit 1
fi

# Count total customers to seed
TOTAL_CUSTOMERS=$(jq '.customers | length' "$MOCKDATA_FILE")
echo "📊 Found $TOTAL_CUSTOMERS customers to seed"

# Counter for progress
SUCCESS_COUNT=0
ERROR_COUNT=0

# Read each customer and create via GraphQL mutation
jq -c '.customers[]' "$MOCKDATA_FILE" | while read -r customer; do
    # Extract customer data
    NAME=$(echo "$customer" | jq -r '.name')
    EMAIL=$(echo "$customer" | jq -r '.email // ""')
    PHONE=$(echo "$customer" | jq -r '.phone // ""')
    SEGMENT=$(echo "$customer" | jq -r '.segment // "RETAIL"')
    STATUS=$(echo "$customer" | jq -r '.status // "ACTIVE"')
    RISK_PROFILE=$(echo "$customer" | jq -r '.riskProfile // "MODERATE"')
    
    # Build traits JSON
    TRAITS_JSON=$(echo "$customer" | jq -c '{
        personality: (.traits.personality // []),
        preferences: (.traits.preferences // []),
        notes: (.traits.notes // ""),
        tags: (.traits.tags // []),
        properties: (.traits.properties // [] | map({key: .key, value: .value, type: .type}))
    }')
    
    # Build GraphQL mutation with proper escaping
    MUTATION=$(cat <<EOF
mutation CreateCustomer(\$input: CreateCustomerInput!) {
  createCustomer(input: \$input) {
    id
    name
    email
  }
}
EOF
)

    VARIABLES=$(cat <<EOF
{
  "input": {
    "name": $(echo "$NAME" | jq -R .),
    "email": $(echo "$EMAIL" | jq -R .),
    "phone": $(echo "$PHONE" | jq -R .),
    "segment": "$SEGMENT",
    "status": "$STATUS",
    "riskProfile": "$RISK_PROFILE",
    "traits": $TRAITS_JSON
  }
}
EOF
)

    # Execute mutation
    RESPONSE=$(curl -s -X POST "$GRAPHQL_URL" \
        -H "Content-Type: application/json" \
        -d "$(jq -n --arg query "$MUTATION" --argjson variables "$VARIABLES" '{query: $query, variables: $variables}')")
    
    # Check for errors
    if echo "$RESPONSE" | jq -e '.errors' > /dev/null 2>&1; then
        ERROR_COUNT=$((ERROR_COUNT + 1))
        echo "  ❌ Failed: $NAME ($EMAIL)"
        echo "     Error: $(echo "$RESPONSE" | jq -r '.errors[0].message')"
    elif echo "$RESPONSE" | jq -e '.data.createCustomer.id' > /dev/null 2>&1; then
        SUCCESS_COUNT=$((SUCCESS_COUNT + 1))
        CREATED_ID=$(echo "$RESPONSE" | jq -r '.data.createCustomer.id')
        echo "  ✅ Created: $NAME ($CREATED_ID)"
    else
        ERROR_COUNT=$((ERROR_COUNT + 1))
        echo "  ❌ Unknown error for: $NAME"
        echo "     Response: $RESPONSE"
    fi
done

echo ""
echo "======================================"
echo "📊 Seeding Summary"
echo "======================================"
echo "  Total attempted: $TOTAL_CUSTOMERS"
echo "  Successfully created: $SUCCESS_COUNT"
echo "  Errors: $ERROR_COUNT"
echo ""

# Verify final count
FINAL_COUNT=$(curl -s -X POST "$GRAPHQL_URL" \
    -H "Content-Type: application/json" \
    -d '{"query":"{ customers { totalCount } }"}' \
    | jq -r '.data.customers.totalCount' 2>/dev/null || echo "0")

echo "✅ Database now has $FINAL_COUNT customers"

if [ "$FINAL_COUNT" -ge "1" ]; then
    echo ""
    echo "🎉 Ready to test collaborative filtering!"
    echo "   Run: ./test-quick.sh"
else
    echo ""
    echo "⚠️  Warning: Database appears empty. Check the errors above."
    exit 1
fi
