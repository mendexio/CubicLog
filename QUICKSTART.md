# CubicLog Quick Start Guide

Get up and running with CubicLog in minutes. This practical guide shows you exactly how to use CubicLog effectively.

## Quick Setup

### 1. Get CubicLog
```bash
# Download binary (replace with your platform)
wget https://github.com/mendexio/CubicLog/releases/latest/download/cubiclog-linux-amd64.tar.gz
tar -xzf cubiclog-linux-amd64.tar.gz
chmod +x cubiclog-linux-amd64
mv cubiclog-linux-amd64 cubiclog

# Or build from source
git clone https://github.com/mendexio/CubicLog.git && cd CubicLog
go build -o cubiclog
```

### 2. Start CubicLog
```bash
./cubiclog
# Server running at http://localhost:8080
```

### 3. Send Your First Log
```bash
curl -X POST http://localhost:8080/api/logs \
  -H 'Content-Type: application/json' \
  -d '{"header": {"title": "Hello CubicLog!"}}'
```

## Basic Logging Patterns

### Minimal Logging (Recommended)
Just send a title - CubicLog figures out the rest:

```bash
# Error log
curl -X POST http://localhost:8080/api/logs \
  -d '{"header": {"title": "Database connection failed"}}'
# → Auto-detected: type=error, color=red

# Success log  
curl -X POST http://localhost:8080/api/logs \
  -d '{"header": {"title": "Payment processed successfully"}}'
# → Auto-detected: type=success, color=green

# Warning log
curl -X POST http://localhost:8080/api/logs \
  -d '{"header": {"title": "High memory usage detected"}}'
# → Auto-detected: type=warning, color=yellow
```

### With Context Data
Add any data you want in the body:

```bash
curl -X POST http://localhost:8080/api/logs \
  -H 'Content-Type: application/json' \
  -d '{
    "header": {
      "title": "User login failed"
    },
    "body": {
      "user_id": 12345,
      "ip": "192.168.1.100",
      "error": "Invalid password",
      "service": "auth-api"
    }
  }'
# CubicLog auto-extracts: source=auth-api, type=error, color=red
```

### Full Control
Specify everything if you prefer:

```bash
curl -X POST http://localhost:8080/api/logs \
  -H 'Content-Type: application/json' \
  -d '{
    "header": {
      "type": "custom",
      "title": "Custom event occurred",
      "description": "Detailed description here",
      "source": "my-service",
      "color": "purple"
    },
    "body": {
      "custom_data": "anything you want"
    }
  }'
```

## Common Use Cases

### Application Errors
```bash
# With stack trace (auto-detected)
curl -X POST http://localhost:8080/api/logs \
  -d '{
    "header": {"title": "Application crashed"},
    "body": {
      "error": "java.lang.NullPointerException at com.example.Service.process(Service.java:142)",
      "service": "payment-processor"
    }
  }'
```

### HTTP Requests
```bash
# API call logging
curl -X POST http://localhost:8080/api/logs \
  -d '{
    "header": {"title": "API request completed"},
    "body": {
      "method": "POST",
      "url": "/api/users",
      "status": 201,
      "duration_ms": 1250,
      "service": "user-api"
    }
  }'
```

### System Monitoring
```bash
# Performance alert
curl -X POST http://localhost:8080/api/logs \
  -d '{
    "header": {"title": "CPU usage high"},
    "body": {
      "cpu_percent": 95,
      "memory_percent": 78,
      "service": "monitoring-agent"
    }
  }'
```

### Business Events
```bash
# Business logic events
curl -X POST http://localhost:8080/api/logs \
  -d '{
    "header": {"title": "Order completed"},
    "body": {
      "order_id": "ORD-12345",
      "amount": 99.99,
      "user_id": 67890,
      "service": "order-service"
    }
  }'
```

## Language Examples

### Node.js/JavaScript
```javascript
const axios = require('axios');

async function log(title, data = {}) {
  try {
    await axios.post('http://localhost:8080/api/logs', {
      header: { title },
      body: data
    });
  } catch (error) {
    console.error('Logging failed:', error.message);
  }
}

// Usage
log('User registered', { user_id: 123, email: 'user@example.com' });
log('Database error occurred', { error: 'Connection timeout' });
```

### Python
```python
import requests
import json

def log(title, data=None):
    try:
        requests.post('http://localhost:8080/api/logs', 
                     json={'header': {'title': title}, 'body': data or {}})
    except Exception as e:
        print(f"Logging failed: {e}")

# Usage
log('Process started', {'process_id': 12345})
log('Error in data processing', {'error': 'File not found', 'file': 'data.csv'})
```

### Go
```go
package main

import (
    "bytes"
    "encoding/json"
    "net/http"
)

func log(title string, data map[string]interface{}) {
    payload := map[string]interface{}{
        "header": map[string]string{"title": title},
        "body": data,
    }
    
    jsonData, _ := json.Marshal(payload)
    http.Post("http://localhost:8080/api/logs", "application/json", 
              bytes.NewBuffer(jsonData))
}

// Usage
log("Service started", map[string]interface{}{"port": 8080})
log("Database connection failed", map[string]interface{}{"error": "timeout"})
```

### Bash/Shell Scripts
```bash
#!/bin/bash
log() {
  curl -s -X POST http://localhost:8080/api/logs \
    -H 'Content-Type: application/json' \
    -d "{\"header\": {\"title\": \"$1\"}, \"body\": {\"message\": \"$2\"}}"
}

# Usage
log "Backup started" "Daily backup process initiated"
log "Disk space low" "Only 5% remaining on /var"
```

## Viewing and Searching Logs

### Web Dashboard
Open http://localhost:8080 for the full dashboard with:
- Real-time log stream
- Advanced filtering
- Smart pattern analytics  
- Export capabilities
- Dark/light mode

### API Access
```bash
# Get all logs
curl http://localhost:8080/api/logs

# Search logs
curl "http://localhost:8080/api/logs?q=database"

# Filter by type
curl "http://localhost:8080/api/logs?type=error"

# Filter by color
curl "http://localhost:8080/api/logs?color=red"

# Date range
curl "http://localhost:8080/api/logs?from=2024-01-01&to=2024-01-31"

# Combine filters
curl "http://localhost:8080/api/logs?type=error&q=timeout&limit=50"
```

### Export Data
```bash
# Export as CSV
curl "http://localhost:8080/api/export/csv" > logs.csv

# Export as JSON
curl "http://localhost:8080/api/export/json" > logs.json

# Export filtered data
curl "http://localhost:8080/api/export/csv?type=error&from=2024-01-01" > errors.csv
```

## Configuration

### Environment Variables
```bash
# Server settings
export PORT=8080
export API_KEY=your-secret-key  # Optional authentication

# Database settings
export DB_PATH=./logs.db
export RETENTION_DAYS=30        # Auto-cleanup after 30 days
```

### Command Line Options
```bash
./cubiclog -port 8081           # Different port
./cubiclog -api-key secret123   # Enable API authentication  
./cubiclog -db /path/logs.db    # Custom database location
./cubiclog -retention 60        # Keep logs for 60 days
./cubiclog -cleanup             # Run cleanup and exit
./cubiclog -version             # Show version
```

### Using with Authentication
```bash
# Set API key
./cubiclog -api-key mysecret

# Send authenticated requests
curl -X POST http://localhost:8080/api/logs \
  -H 'Authorization: Bearer mysecret' \
  -H 'Content-Type: application/json' \
  -d '{"header": {"title": "Authenticated log"}}'
```

## Smart Pattern Detection

CubicLog automatically detects and categorizes logs:

### HTTP Status Codes
```bash
# These get auto-detected and categorized
echo '{"header": {"title": "Request returned 404"}}' # → type=warning, color=yellow
echo '{"header": {"title": "Server error 500"}}' # → type=error, color=red  
echo '{"header": {"title": "Request successful 200"}}' # → type=success, color=green
```

### Stack Traces
```bash
# Auto-detected as errors
echo '{
  "header": {"title": "Application error"},
  "body": {"error": "java.lang.NullPointerException at com.example.Service.process"}
}' # → Detected stack trace, severity=error
```

### Performance Issues
```bash
# Auto-detected performance problems
echo '{
  "header": {"title": "Slow database query"},
  "body": {"duration_ms": 5500}
}' # → Detected slow query, severity=warning
```

## Troubleshooting

### Common Issues

**Port already in use:**
```bash
./cubiclog -port 8081
```

**Permission denied:**
```bash
chmod +x cubiclog
```

**Can't connect to CubicLog:**
```bash
# Check if running
curl http://localhost:8080/health
# Should return: {"status":"ok"}
```

**Database issues:**
```bash
# Reset database
rm logs.db
./cubiclog  # Will recreate automatically
```

**High disk usage:**
```bash
# Manual cleanup
./cubiclog -cleanup

# Or reduce retention
./cubiclog -retention 7  # Keep only 7 days
```

### Getting Help

**Check logs are being received:**
```bash
# Send test log
curl -X POST http://localhost:8080/api/logs \
  -d '{"header": {"title": "Test log"}}'

# Check it appears
curl http://localhost:8080/api/logs | jq '.[0]'
```

**Check smart pattern detection:**
```bash
# Get pattern statistics
curl http://localhost:8080/api/stats | jq '.pattern_stats'

# Example output:
# {
#   "http_codes_detected": 45,
#   "stack_traces_found": 12, 
#   "security_issues": 3,
#   "performance_issues": 28
# }
```

**View detection accuracy:**
```bash
curl http://localhost:8080/api/stats | jq '.detection_accuracy'
# "91.5%"
```

## Tips for Effective Logging

### 1. Start Simple
Just send a title - let CubicLog handle the categorization:
```bash
curl -X POST http://localhost:8080/api/logs \
  -d '{"header": {"title": "Something interesting happened"}}'
```

### 2. Include Service Context
Help CubicLog identify sources by including service info:
```bash
curl -X POST http://localhost:8080/api/logs \
  -d '{
    "header": {"title": "User updated profile"},
    "body": {"service": "user-service", "user_id": 123}
  }'
```

### 3. Use Structured Data
Put variable data in the body for better analysis:
```bash
curl -X POST http://localhost:8080/api/logs \
  -d '{
    "header": {"title": "Payment processed"},
    "body": {"amount": 99.99, "currency": "USD", "gateway": "stripe"}
  }'
```

### 4. Let Patterns Work
Include error details for automatic severity detection:
```bash
curl -X POST http://localhost:8080/api/logs \
  -d '{
    "header": {"title": "Operation failed"},
    "body": {"error": "Connection timeout after 30s"}
  }'
```

---

That's it! You're now ready to use CubicLog effectively. 

**Next Steps:**
- 📖 [Full Documentation](README.md) - Comprehensive features and deployment guide
- 🌐 [GitHub Pages](https://mendexio.github.io/CubicLog/) - Interactive documentation and examples
- 🐛 [Issues & Support](https://github.com/mendexio/CubicLog/issues) - Get help or report problems