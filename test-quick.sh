#!/bin/bash
# test-quick.sh - Quick test script using Makefile commands
# Usage: ./test-quick.sh

echo "🧪 Quick Collaborative Filtering Test"
echo "======================================"
echo ""

# Detect GraphQL port (try both 8080 and 4000)
if curl -s http://localhost:8080/health > /dev/null 2>&1 || curl -s -X POST http://localhost:8080/query -d '{"query":"{ __typename }"}' > /dev/null 2>&1; then
    BASE_URL="http://localhost:8080/query"
    PORT="8080"
elif curl -s http://localhost:4000/health > /dev/null 2>&1 || curl -s -X POST http://localhost:4000/query -d '{"query":"{ __typename }"}' > /dev/null 2>&1; then
    BASE_URL="http://localhost:4000/query"
    PORT="4000"
else
    echo "❌ GraphQL server not responding on port 8080 or 4000"
    echo ""
    echo "Please start services first:"
    echo "  make dev"
    echo ""
    exit 1
fi

echo "✅ GraphQL server detected on port $PORT"

# Check if services are running
if ! docker-compose -f docker-compose.dev.yml ps | grep -q "Up"; then
    echo "⚠️  Services not running. Starting now..."
    make dev
    echo "⏳ Waiting 10 seconds for services to be ready..."
    sleep 10
fi

# Check if data is loaded
echo ""
echo "📊 Checking database..."
CUSTOMER_COUNT=$(curl -s -X POST $BASE_URL \
  -H "Content-Type: application/json" \
  -d '{"query":"{ customers { totalCount } }"}' \
  | jq -r '.data.customers.totalCount' 2>/dev/null)

# Handle empty or null response
if [ -z "$CUSTOMER_COUNT" ] || [ "$CUSTOMER_COUNT" = "null" ] || [ "$CUSTOMER_COUNT" = "0" ]; then
    echo "⚠️  No customers in database (found: ${CUSTOMER_COUNT:-0})"
    echo ""
    echo "Please seed the database first:"
    echo "  ./seed-data.sh"
    echo ""
    echo "Or use the full setup:"
    echo "  make test-setup"
    exit 1
fi

echo "✅ Found ${CUSTOMER_COUNT:-0} customers in database"

echo ""
echo "🧪 Running Critical Tests..."
echo ""

# Get first customer ID from database
CUSTOMER_ID=$(curl -s -X POST $BASE_URL \
  -H "Content-Type: application/json" \
  -d '{"query":"{ customers(limit: 1) { nodes { id } } }"}' \
  | jq -r '.data.customers.nodes[0].id' 2>/dev/null)

if [ -z "$CUSTOMER_ID" ] || [ "$CUSTOMER_ID" = "null" ]; then
    echo "❌ Could not retrieve customer ID from database"
    exit 1
fi

echo "📝 Using customer ID: $CUSTOMER_ID"
echo ""

# Test 1: Similar Customers
echo "Test 1: Find Similar Customers"
SIMILAR_RESPONSE=$(curl -s -X POST $BASE_URL \
  -H "Content-Type: application/json" \
  -d "{\"query\":\"{ similarCustomers(customerId: \\\"$CUSTOMER_ID\\\", limit: 3) { customer { name } similarityScore } }\"}")

SIMILAR_COUNT=$(echo "$SIMILAR_RESPONSE" | jq -r '.data.similarCustomers | length' 2>/dev/null)

if [ -n "$SIMILAR_COUNT" ] && [ "$SIMILAR_COUNT" != "null" ] && [ "$SIMILAR_COUNT" -ge "1" ] 2>/dev/null; then
    echo "  ✅ PASSED - Found $SIMILAR_COUNT similar customers"
else
    echo "  ❌ FAILED - No similar customers found"
    echo "     Response: $(echo "$SIMILAR_RESPONSE" | jq -c '.' 2>/dev/null || echo "$SIMILAR_RESPONSE")"
fi

# Test 2: Generate Recommendations
echo "Test 2: Generate Recommendations"
REC_RESPONSE=$(curl -s -X POST $BASE_URL \
  -H "Content-Type: application/json" \
  -d "{\"query\":\"{ generateRecommendations(customerId: \\\"$CUSTOMER_ID\\\", type: FUND) { itemName } }\"}")

REC_COUNT=$(echo "$REC_RESPONSE" | jq -r '.data.generateRecommendations | length' 2>/dev/null)

if [ -n "$REC_COUNT" ] && [ "$REC_COUNT" != "null" ] && [ "$REC_COUNT" -ge "1" ] 2>/dev/null; then
    echo "  ✅ PASSED - Generated $REC_COUNT recommendations"
else
    echo "  ❌ FAILED - No recommendations generated"
    echo "     Response: $(echo "$REC_RESPONSE" | jq -c '.' 2>/dev/null || echo "$REC_RESPONSE")"
fi

# Test 3: Get Customer
echo "Test 3: Get Customer Profile"
CUSTOMER_RESPONSE=$(curl -s -X POST $BASE_URL \
  -H "Content-Type: application/json" \
  -d "{\"query\":\"{ customer(id: \\\"$CUSTOMER_ID\\\") { name } }\"}")

CUSTOMER_NAME=$(echo "$CUSTOMER_RESPONSE" | jq -r '.data.customer.name' 2>/dev/null)

if [ -n "$CUSTOMER_NAME" ] && [ "$CUSTOMER_NAME" != "null" ]; then
    echo "  ✅ PASSED - Retrieved customer: $CUSTOMER_NAME"
else
    echo "  ❌ FAILED - Could not retrieve customer"
    echo "     Response: $(echo "$CUSTOMER_RESPONSE" | jq -c '.' 2>/dev/null || echo "$CUSTOMER_RESPONSE")"
fi

echo ""
echo "======================================"
echo "🎉 Quick test complete!"
echo ""
echo "Next steps:"
echo "  - View logs: make dev-logs"
echo "  - Stop services: make dev-stop"
echo "  - Full tests: ./test-all.sh"
echo "  - GraphQL Playground: http://localhost:8080"
