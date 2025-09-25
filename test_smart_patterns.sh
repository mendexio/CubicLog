#!/bin/bash
# Test script for CubicLog Smart Pattern Matching

echo "🧪 Testing CubicLog Smart Pattern Matching"
echo "=========================================="

# Start CubicLog if not running
if ! pgrep -f cubiclog > /dev/null; then
    echo "Starting CubicLog..."
    ./cubiclog &
    sleep 2
fi

URL="http://localhost:8080"

echo "\n1️⃣ Testing HTTP Status Codes"
curl -X POST $URL/api/logs -H "Content-Type: application/json" -d '{
    "header": {"title": "API Response"},
    "body": {"message": "Request returned 500 Internal Server Error"}
}'

echo "\n2️⃣ Testing Stack Trace Detection"
curl -X POST $URL/api/logs -H "Content-Type: application/json" -d '{
    "header": {"title": "Application Error"},
    "body": {
        "error": "java.lang.NullPointerException at com.example.Service.process(Service.java:142)"
    }
}'

echo "\n3️⃣ Testing Security Pattern"
curl -X POST $URL/api/logs -H "Content-Type: application/json" -d '{
    "header": {"title": "Security Alert"},
    "body": {"message": "Unauthorized access attempt detected, possible SQL injection"}
}'

echo "\n4️⃣ Testing Performance Metrics"
curl -X POST $URL/api/logs -H "Content-Type: application/json" -d '{
    "header": {"title": "Slow Query"},
    "body": {"message": "Database query took 5234ms to complete"}
}'

echo "\n5️⃣ Testing Database Error"
curl -X POST $URL/api/logs -H "Content-Type: application/json" -d '{
    "header": {"title": "DB Error"},
    "body": {"error": "Deadlock detected when acquiring lock on users table"}
}'

echo "\n6️⃣ Testing System Error Code"
curl -X POST $URL/api/logs -H "Content-Type: application/json" -d '{
    "header": {"title": "Connection Failed"},
    "body": {"error": "Failed to connect: ECONNREFUSED"}
}'

echo "\n7️⃣ Testing Business Logic"
curl -X POST $URL/api/logs -H "Content-Type: application/json" -d '{
    "header": {"title": "Payment Update"},
    "body": {"message": "Payment successful for order #12345", "amount": 99.99}
}'

echo "\n8️⃣ Testing Resource Usage"
curl -X POST $URL/api/logs -H "Content-Type: application/json" -d '{
    "header": {"title": "System Monitor"},
    "body": {"message": "CPU: 95%, Memory: 88%, Disk: 45%"}
}'

echo "\n\n📊 Checking Smart Analytics..."
sleep 2
curl -s $URL/api/stats | python -m json.tool | grep -E '"(pattern_stats|detection_accuracy|severity_breakdown)"' -A 5

echo "\n✅ Smart Pattern Testing Complete!"
echo "Check dashboard at $URL to see categorized logs"