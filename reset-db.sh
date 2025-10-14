#!/bin/bash

# Reset database by restarting Docker containers with fresh volumes
# This is the simplest way to get a clean database

echo "🗑️  Resetting Database..."
echo "======================================"
echo ""
echo "⚠️  WARNING: This will delete ALL database data!"
echo ""
read -p "Are you sure you want to continue? (yes/no): " -r
echo ""

if [[ ! $REPLY =~ ^[Yy]es$ ]]; then
    echo "❌ Reset cancelled"
    exit 0
fi

echo "🗑️  Step 1: Stopping containers..."
docker-compose -f docker-compose.dev.yml down -v

echo ""
echo "🚀 Step 2: Starting containers with fresh database..."
docker-compose -f docker-compose.dev.yml up -d

echo ""
echo "⏳ Step 3: Waiting for services to start (15 seconds)..."
sleep 15

echo ""
echo "======================================"
echo "✅ Database Reset Complete!"
echo "======================================"
echo ""
echo "Next steps:"
echo "  1. Check if server is running: curl http://localhost:8080/health"
echo "  2. Load mock data: ./seed-data.sh"
echo "  3. Run tests: ./test-quick.sh"
echo ""
echo "If server isn't ready, check logs:"
echo "  make dev-logs"
