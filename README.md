# CubicLog

A beautifully simple self-hosted logging solution. No Kubernetes, no complexity, just logs.

![CubicLog Dashboard](https://img.shields.io/badge/CubicLog-v1.0.0-blue) ![Go](https://img.shields.io/badge/Go-1.21+-blue) ![SQLite](https://img.shields.io/badge/SQLite-embedded-green)

## Features

### 🧠 **Intelligent Analytics**
- 🤖 **Automated metadata derivation** - Smart extraction of severity, sources, and categories
- 📈 **Real-time error rate monitoring** with trend analysis
- 🚨 **Smart alerts system** - Automatic notifications for high error rates and anomalies
- 📊 **Server health dashboard** - Color-coded system status indicators
- 🔄 **Volume trend analysis** - Activity level monitoring (increasing/stable/decreasing)
- 🎯 **Pattern recognition** - Intelligent severity detection from log content

### 🎨 **Beautiful Interface**
- 🎨 **22 Tailwind CSS colors** for comprehensive log categorization
- 🔍 **Advanced search & filtering** with real-time results
- 🌙 **Dark/light mode** toggle with persistence
- 📱 **Responsive design** - Works perfectly on mobile and desktop
- ⚡ **Real-time updates** - 5-second auto-refresh with live data
- 🎭 **JSON syntax highlighting** with collapsible structures

### 🏗️ **Core Infrastructure**
- 📁 **SQLite storage** - No external database required
- 🔐 **API key authentication** - Optional security layer
- 📦 **Single binary deployment** - Download and run
- 📤 **CSV/JSON export** - Filtered data export capabilities
- 🧹 **Automatic log retention** - Configurable cleanup policies
- 🛠️ **Service management** - Start/stop/restart/status commands

### ⚡ **Performance & Reliability**
- 🚀 **Lightning fast** - Runs efficiently on Raspberry Pi
- 🔄 **Auto-recovery** - Graceful error handling and restart capabilities
- 💾 **Minimal footprint** - ~11MB single binary with embedded UI
- 🌐 **Cross-platform** - Linux, Windows, macOS (Intel & Apple Silicon)

## Quick Start

### Install

**Option 1: Download Binary**
```bash
# Download from GitHub Releases
wget https://github.com/mendexio/CubicLog/releases/latest/download/cubiclog-linux-amd64.tar.gz
tar -xzf cubiclog-linux-amd64.tar.gz
chmod +x cubiclog-linux-amd64
mv cubiclog-linux-amd64 cubiclog
```

**Option 2: Build from Source**
```bash
git clone https://github.com/mendexio/CubicLog.git
cd CubicLog
go build -o cubiclog
```

**Option 3: Go Install**
```bash
go install github.com/mendexio/CubicLog@latest
```

### Run
```bash
./cubiclog
```

### Visit
Open [http://localhost:8080](http://localhost:8080) in your browser.

That's it! 🎉

## Usage

### Send a Log

**⚠️ CubicLog v1.0 requires ALL header fields to be provided:**

**Complete Log Entry:**
```bash
curl -X POST http://localhost:8080/api/logs \
  -H 'Content-Type: application/json' \
  -d '{
    "header": {
      "type": "info",
      "title": "User logged in",
      "description": "Successful user authentication from web interface",
      "source": "auth-service",
      "color": "blue"
    },
    "body": {
      "user_id": 123,
      "ip": "192.168.1.1",
      "user_agent": "Mozilla/5.0...",
      "session_id": "abc123"
    }
  }'
```

**Error Log Example:**
```bash
curl -X POST http://localhost:8080/api/logs \
  -H 'Content-Type: application/json' \
  -d '{
    "header": {
      "type": "error",
      "title": "Payment Processing Failed",
      "description": "Credit card transaction was declined by payment gateway",
      "source": "payment-service",
      "color": "red"
    },
    "body": {
      "user_id": 123,
      "amount": 99.99,
      "error_code": "CARD_DECLINED",
      "transaction_id": "txn_abc123",
      "gateway_response": "Insufficient funds"
    }
  }'
```

### Required Header Fields

**All fields are mandatory in v1.0:**

| Field | Type | Description | Example |
|-------|------|-------------|---------|
| `type` | string | Log category/type | `"error"`, `"info"`, `"warning"` |
| `title` | string | Brief, descriptive title | `"Payment Failed"` |
| `description` | string | Detailed explanation | `"Credit card was declined by gateway"` |
| `source` | string | Originating service/component | `"payment-service"`, `"auth-api"` |
| `color` | string | Valid Tailwind CSS 4 color | `"red"`, `"blue"`, `"green"` |

### Color Options

CubicLog supports all Tailwind CSS v4 named colors for visual organization:

**Primary Colors:**
- <span style="color: #ef4444">**red**</span> - Errors, failures, critical issues
- <span style="color: #22c55e">**green**</span> - Success, completed operations  
- <span style="color: #eab308">**yellow**</span> - Warnings, alerts
- <span style="color: #3b82f6">**blue**</span> - Info, general logs
- <span style="color: #a855f7">**purple**</span> - Custom events
- <span style="color: #6b7280">**gray**</span> - Debug, trace logs

**Extended Palette:**
- <span style="color: #f97316">**orange**</span> - Performance issues, slow operations
- <span style="color: #ec4899">**pink**</span> - User interactions, UI events
- <span style="color: #6366f1">**indigo**</span> - Database operations
- <span style="color: #06b6d4">**cyan**</span> - Network requests, API calls
- <span style="color: #64748b">**slate**</span> - System logs, background tasks
- <span style="color: #71717a">**zinc**</span> - Cache operations
- <span style="color: #737373">**neutral**</span> - Default/unclassified logs
- <span style="color: #78716c">**stone**</span> - File operations

**Bright Colors:**
- <span style="color: #65a30d">**lime**</span> - High priority success
- <span style="color: #059669">**emerald**</span> - Completion events
- <span style="color: #0d9488">**teal**</span> - Data processing
- <span style="color: #0ea5e9">**sky**</span> - Cloud/external services
- <span style="color: #8b5cf6">**violet**</span> - Authentication events
- <span style="color: #d946ef">**fuchsia**</span> - Special notifications
- <span style="color: #f43f5e">**rose**</span> - Security alerts

**⚠️ Color Required**: In v1.0, the `color` field is mandatory and must be one of the 22 valid Tailwind CSS 4 color names listed above. Auto-assignment has been removed to ensure consistent, intentional color usage.

## 🧠 Intelligent Analytics

CubicLog v1.0 introduces powerful intelligent analytics that automatically derive insights from your logs without requiring structured schemas or complex configuration.

### Philosophy: "Be Liberal in What You Accept, Intelligent in What You Derive"

CubicLog allows you to send **any JSON structure** in the log body while automatically extracting meaningful insights through pattern recognition and intelligent analysis.

### 🤖 Automated Metadata Derivation

Every log automatically gets analyzed to derive:

**Severity Detection:**
- **Error**: Detects keywords like "error", "failed", "exception", "crash", "critical"
- **Warning**: Identifies "warning", "slow", "timeout", "deprecated", "retry"  
- **Success**: Recognizes "success", "completed", "approved", "validated"
- **Debug**: Finds "debug", "trace", "verbose", "entering", "exiting"
- **Info**: Default fallback for unmatched patterns

**Smart Source Extraction:**
- Automatically extracts service names from `body.service`, `body.source`, or `header.source`
- Intelligently identifies the most relevant source identifier
- Powers the analytics dashboard with accurate service mapping

**Category Classification:**
- Derives categories from log types or meaningful title keywords
- Enables automatic log organization and filtering

### 📊 Real-Time Analytics Dashboard

The intelligent dashboard provides instant insights:

**Server Health Monitoring:**
- 🟢 **Healthy** (< 10% error rate): "All systems go"
- 🟡 **Warning** (10-30% error rate): "Monitor closely"  
- 🔴 **Critical** (> 30% error rate): "Needs attention"

**Smart Alerts:**
- Automatically appear for high error rates
- Color-coded severity levels (high/medium/low)
- Actionable insights with automated detection

**Trend Analysis:**
- **Volume Trends**: Increasing ↑ / Stable = / Decreasing ↓
- **Error Rate Trending**: Real-time calculation and historical analysis
- **Activity Monitoring**: Live system health indicators

**Visual Analytics:**
- Error rate percentages with trend indicators
- Log type distribution with interactive charts
- Source analysis showing top log generators
- Hourly distribution patterns for anomaly detection

### 🔍 Analytics API

Access analytics programmatically:

```bash
# Get comprehensive analytics
curl http://localhost:8080/api/stats

# Example response
{
  "total": 1234,
  "error_rate_24h": "15.5%",
  "severity_breakdown": {
    "error": 45,
    "warning": 23,
    "success": 156,
    "info": 891,
    "debug": 119
  },
  "top_sources": [
    {"name": "payment-gateway", "count": 234},
    {"name": "user-service", "count": 189}
  ],
  "alerts": ["Error rate above 15% - monitor closely"],
  "trends": {
    "error_trend": "increasing",
    "volume_trend": "stable"
  }
}
```

### 💡 Intelligent Insights Examples

**Automatic Error Detection:**
```json
{
  "header": {
    "type": "payment_issue",
    "title": "Transaction failed",
    "description": "Credit card processing failed with timeout",
    "source": "billing",
    "color": "red"
  },
  "body": {
    "service": "stripe-integration",
    "amount": 99.99,
    "error": "Gateway timeout"
  }
}
```
**→ Automatically derived: severity="error", source="stripe-integration"**

**Smart Success Recognition:**
```json
{
  "header": {
    "type": "user_action",
    "title": "Login completed",
    "description": "User authentication successful",
    "source": "auth",
    "color": "green"
  },
  "body": {
    "service": "user-authentication-api",
    "user_id": 12345
  }
}
```
**→ Automatically derived: severity="success", source="user-authentication-api"**

The intelligent analytics work **behind the scenes** with zero configuration, turning your unstructured logs into actionable insights automatically.

## Configuration

All configuration is optional via environment variables:

```bash
# Server settings
PORT=8080                    # Port to run on (default: 8080)
API_KEY=your-secret-key     # Enable authentication (optional)

# Database settings  
DB_PATH=./logs.db           # SQLite database path (default: ./logs.db)
RETENTION_DAYS=30           # Days to keep logs (default: 30)
```

### CLI Flags

```bash
./cubiclog --help

Usage of ./cubiclog:
  -api-key string
        API key for authentication
  -cleanup
        Run cleanup and exit
  -db string
        Path to SQLite database (default "./logs.db")
  -port string
        Port to run server on (default "8080")
  -retention int
        Days to retain logs (default 30)
  -version
        Show version
```

## API Endpoints

### Logs
- `POST /api/logs` - Send logs
- `GET /api/logs` - View logs (supports filters)
- `GET /api/stats` - Statistics
- `GET /health` - Health check

### Export  
- `GET /api/export/csv` - Export as CSV
- `GET /api/export/json` - Export as JSON

### Filters

All endpoints support these query parameters:

| Parameter | Description | Example |
|-----------|-------------|---------|
| `q` | Search text | `?q=database` |
| `type` | Log type | `?type=error` |
| `color` | Log color | `?color=red` |
| `from` | Start date | `?from=2024-01-01` |
| `to` | End date | `?to=2024-01-31` |
| `limit` | Max results | `?limit=50` |

**Combine filters:**
```bash
GET /api/logs?type=error&color=red&q=timeout&limit=50
```

## Examples

### Application Logging

**Node.js:**
```javascript
const axios = require('axios');

function log(type, title, description, color, data = {}) {
  axios.post('http://localhost:8080/api/logs', {
    header: { 
      type, 
      title, 
      description, 
      source: 'my-app',
      color 
    },
    body: data
  }).catch(console.error);
}

// Usage
log('error', 'Database connection failed', 
    'Failed to establish connection to PostgreSQL database',
    'red', 
    { 
      error: 'ECONNREFUSED', 
      host: 'localhost:5432',
      timeout: 5000
    });
```

**Python:**
```python
import requests

def log(log_type, title, description, source, color, data=None):
    requests.post('http://localhost:8080/api/logs', json={
        'header': {
            'type': log_type, 
            'title': title, 
            'description': description,
            'source': source,
            'color': color
        },
        'body': data or {}
    })

# Usage  
log('info', 'Process completed', 
    'Data processing pipeline finished successfully',
    'data-processor', 'green',
    {'processed': 1000, 'errors': 0, 'duration_ms': 5420})
```

**Go:**
```go
type LogEntry struct {
    Header struct {
        Type        string `json:"type"`
        Title       string `json:"title"`
        Description string `json:"description"`
        Source      string `json:"source"`
        Color       string `json:"color"`
    } `json:"header"`
    Body map[string]interface{} `json:"body"`
}

func sendLog(logType, title, description, source, color string, data map[string]interface{}) {
    entry := LogEntry{
        Header: struct {
            Type        string `json:"type"`
            Title       string `json:"title"`
            Description string `json:"description"`
            Source      string `json:"source"`
            Color       string `json:"color"`
        }{
            Type: logType, 
            Title: title, 
            Description: description,
            Source: source, 
            Color: color,
        },
        Body: data,
    }
    // Send to CubicLog via HTTP POST...
}

// Usage
sendLog("error", "Redis Connection Failed", 
        "Failed to connect to Redis cache server",
        "cache-service", "red",
        map[string]interface{}{
            "host": "redis.internal:6379",
            "timeout_ms": 5000,
            "retry_count": 3,
        })
```

### Server Access Logs

```bash
# Nginx log forwarding
tail -f /var/log/nginx/access.log | while read line; do
  curl -X POST http://localhost:8080/api/logs \
    -H 'Content-Type: application/json' \
    -d "{\"header\":{\"type\":\"info\",\"title\":\"HTTP Request\",\"source\":\"nginx\"},\"body\":{\"log\":\"$line\"}}"
done
```

### System Monitoring

```bash
#!/bin/bash
# Monitor disk space
USAGE=$(df / | tail -1 | awk '{print $5}' | sed 's/%//')
if [ $USAGE -gt 80 ]; then
  curl -X POST http://localhost:8080/api/logs \
    -H 'Content-Type: application/json' \
    -d "{\"header\":{\"type\":\"warning\",\"title\":\"High disk usage\",\"color\":\"yellow\"},\"body\":{\"usage\":\"$USAGE%\",\"mount\":\"/\"}}"
fi
```

## Deployment

### Systemd Service

```bash
# Create service file
sudo tee /etc/systemd/system/cubiclog.service > /dev/null <<EOF
[Unit]
Description=CubicLog
After=network.target

[Service]
Type=simple
User=cubiclog
WorkingDirectory=/opt/cubiclog
ExecStart=/opt/cubiclog/cubiclog
Restart=always
Environment=PORT=8080
Environment=RETENTION_DAYS=30

[Install]
WantedBy=multi-user.target
EOF

# Enable and start
sudo systemctl enable cubiclog
sudo systemctl start cubiclog
```

### Docker (Optional)

```dockerfile
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY cubiclog .
CMD ["./cubiclog"]
```

```bash
docker build -t cubiclog .
docker run -p 8080:8080 -v $(pwd)/data:/root cubiclog
```

### Reverse Proxy

**Nginx:**
```nginx
server {
    listen 80;
    server_name logs.example.com;
    
    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

**Caddy:**
```
logs.example.com {
    reverse_proxy localhost:8080
}
```

## Why CubicLog?

### ✅ **Simple**
- 3 files, 1 dependency
- No build process needed
- Zero configuration required

### ⚡ **Fast** 
- Runs on a Raspberry Pi
- SQLite = no network overhead
- Embedded web UI

### 🎨 **Beautiful**
- Modern, responsive design
- Dark mode support
- Color-coded for quick scanning

### 🔒 **Reliable**
- Your data stays yours
- No external dependencies
- Works offline

### 📦 **Portable**
- Single binary
- Copy and run anywhere
- No installation required

## Architecture

```
┌─────────────────┐    ┌──────────────┐    ┌─────────────┐
│   Web Browser   │◄──►│   CubicLog   │◄──►│   SQLite    │
│  (Dashboard)    │    │   (Go)       │    │ (logs.db)   │
└─────────────────┘    └──────────────┘    └─────────────┘
                              ▲
                              │
                       ┌──────────────┐
                       │ Applications │
                       │   (REST)     │
                       └──────────────┘
```

**File Structure:**
```
cubiclog/
├── main.go      # Core server logic (500 lines)
├── web.go       # Embedded web UI (HTML/CSS/JS)
├── main_test.go # Simple tests
├── README.md    # This file
├── go.mod       # Single dependency: sqlite3
└── logs.db      # Created automatically
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

**Database locked:**
```bash
# Usually resolves itself, or restart CubicLog
sudo systemctl restart cubiclog
```

### Cleanup

**Manual cleanup:**
```bash
./cubiclog -cleanup
```

**Reset database:**
```bash
rm logs.db
./cubiclog  # Will recreate automatically
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Keep it simple - this is the CubicLog way
4. Commit your changes (`git commit -am 'Add amazing feature'`)
5. Push to the branch (`git push origin feature/amazing-feature`)
6. Open a Pull Request

### Development

```bash
# Run tests
go test -v

# Run with hot reload
go run main.go web.go

# Build
go build -o cubiclog
```

## FAQ

**Q: Why not use ELK stack / Splunk / DataDog?**  
A: Sometimes you just want simple logging without the complexity and cost.

**Q: Can it handle high volume?**  
A: SQLite is surprisingly fast. For very high volume, consider log aggregation before sending to CubicLog.

**Q: Is it production ready?**  
A: Yes! It's designed for simplicity and reliability. Many teams use it in production.

**Q: Can I contribute?**  
A: Absolutely! Just remember: simple > complex.

## Development Philosophy

CubicLog was developed following a philosophy of radical simplicity, not only in its final product but also in its creation process. The vision, architecture, and fundamental principles of the project were defined by the human author.

The code implementation was largely accelerated through collaboration with the Claude Code AI assistant. The workflow consisted of providing detailed and incremental engineering instructions, with the author guiding, reviewing, and refining the generated code to ensure it perfectly aligned with the project's philosophy.

This project is an example of how Developer-AI collaboration can be used to build useful and well-designed tools extremely efficiently.

## License

MIT License - Use it however you want.

## Support

- 🐛 **Issues**: [GitHub Issues](https://github.com/mendexio/CubicLog/issues)
- 💬 **Discussions**: [GitHub Discussions](https://github.com/mendexio/CubicLog/discussions)
- 📧 **Email**: support@mendex.io

---

**Built with ❤️ by [Mendex](https://mendex.io)**

*CubicLog - Logging for humans who just want things to work.*