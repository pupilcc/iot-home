# IoT Home - IoT Home Automation Service

A Go-based IoT home automation backend service that connects IoT devices with mobile push notifications. The service receives SMS messages from ESP32 devices via MQTT protocol and forwards them to users' mobile devices through the Bark push notification service.

## Workflow

1. ESP32 device sends SMS data to Broker via MQTT
2. Application subscribes to `esp32/sms` topic and receives JSON messages
3. Parses message content (sender, content, operator, timestamp)
4. Formats Bark notification (title, body, metadata)
5. Sends push notification to user's mobile device via Bark API

## Features

- **MQTT Message Subscription** - Real-time monitoring of SMS data sent by ESP32 devices
- **Push Notification Integration** - Forward messages to mobile devices via Bark service

## Tech Stack

- **Language**: Go 1.24.0
- **Web Framework**: Echo v4.11.3
- **MQTT Client**: Eclipse Paho MQTT v1.5.1
- **Logging**: Uber Zap v1.26.0
- **Containerization**: Docker (Debian Bookworm-slim)
- **CI/CD**: GitHub Actions

## Getting Started

### Prerequisites

- Go 1.24.0 or higher
- Docker (optional)
- MQTT Broker (e.g., Mosquitto)
- Bark push service

### Local Development

1. **Clone the repository**

```bash
git clone <repository-url>
cd iot-home
```

2. **Configure environment variables**

Create a `.env` file and configure the following parameters:

```env
BARK_API=https://xxxx.com
BARK_KEY=your_bark_key_here
MQTT_HOST=192.168.1.1
MQTT_PORT=1883
```

3. **Install dependencies**

```bash
go mod download
```

4. **Run the service**

```bash
go run main.go
```

The service will start at `http://localhost:1323`.

## Configuration

### Environment Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `BARK_API` | Bark push service API URL | `https://xxxx.com` |
| `BARK_KEY` | Bark service authentication key | `xxxx` |
| `MQTT_HOST` | MQTT Broker host address | `192.168.1.1` |
| `MQTT_PORT` | MQTT Broker port | `1883` |

### MQTT Configuration

Default subscription topic: `esp32/sms`

Message format (JSON):
```json
{
  "sender": "13800138000",
  "content": "SMS content",
  "operator": "China Mobile",
  "timestamp": "2025-11-12 10:30:45"
}
```

### Bark Push Configuration

Push notifications support the following features:
- Custom title and content
- Message grouping (Group)
- Notification level (Level): active, timeSensitive, passive
- One-click copy (Copy)
- Markdown format support


## Project Structure

```
iot-home/
├── api/                      # HTTP API routes and handlers
│   └── index.go             # Base endpoint definitions
├── config/                   # Configuration and initialization
│   └── logger.go            # Logger configuration
├── infrastructure/          # Core business logic
│   ├── mqtt/               # MQTT protocol handling
│   │   ├── config.go       # MQTT configuration structure
│   │   ├── connection.go   # MQTT client management
│   │   ├── consumer.go     # Message subscription handling
│   │   └── message.go      # Message structure definitions
│   ├── bark/               # Bark push service
│   │   ├── config.go       # Bark configuration
│   │   └── send.go         # Push logic
│   └── util/               # Utility functions
│       └── env.go          # Environment variable helpers
├── main.go                  # Application entry point
├── go.mod                   # Go module dependencies
├── go.sum                   # Dependency checksums
├── Dockerfile               # Docker build configuration
├── .env                     # Environment variables (development)
└── .github/
    └── workflows/
        └── release-docker-image.yml  # CI/CD pipeline
```

## Security Considerations

- Do not commit `.env` file to version control
- Use strong keys to protect Bark API Key
- Use MQTT authentication in production environments
- Recommended to use HTTPS/TLS encrypted communication

