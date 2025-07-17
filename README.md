# URL Shortener Microservices

Production-ready URL shortener built with Go microservices architecture.

## 🏗️ Architecture

The project consists of three main microservices:
- **URL Service**: Core URL shortening functionality
- **User Service**: Authentication, authorization, and rate limiting
- **Analytics Service**: Click tracking and reporting

### Technology Stack

- **Language**: Go 1.22+
- **Databases**: PostgreSQL, Redis, MongoDB
- **Message Brokers**: NATS, RabbitMQ, Kafka (phased approach)
- **Protocols**: REST, gRPC, WebSocket
- **Observability**: OpenTelemetry, Prometheus, Grafana, Jaeger

## 🚀 Quick Start

### Prerequisites

- Go 1.22+
- Docker & Docker Compose
- Make (optional)

### Development Setup

1. Clone the repository:
```bash
git clone https://github.com/dlaleks/url-shortener-microservices.git
cd url-shortener-microservices
```

2. Initialize the project:
```bash
make init
```

3. Start the development environment:
```bash
make dev-up
```

4. Run the services:
```bash
make build
make run
```

## 📁 Project Structure

```
.
├── services/              # Microservices
│   ├── url-service/      # URL shortening service
│   ├── user-service/     # User management service
│   └── analytics-service/# Analytics service
├── pkg/                  # Shared packages
├── proto/               # Protocol buffer definitions
├── deployments/         # Deployment configurations
├── scripts/            # Utility scripts
└── docs/              # Documentation
```

## 🔧 Development

### Available Commands

```bash
make build              # Build all services
make test              # Run tests
make lint              # Run linters
make dev-up            # Start development environment
make dev-down          # Stop development environment
make logs              # View logs
```

### Service Ports

| Service | Port | Description |
|---------|------|-------------|
| URL Service (HTTP) | 8080 | REST API |
| URL Service (gRPC) | 8081 | gRPC API |
| User Service (HTTP) | 8082 | REST API |
| User Service (gRPC) | 8083 | gRPC API |
| Analytics Service (HTTP) | 8084 | REST API |
| Analytics Service (gRPC) | 8085 | gRPC API |
| PostgreSQL | 5432 | Database |
| Redis | 6379 | Cache |
| MongoDB | 27017 | Analytics DB |
| NATS | 4222 | Message broker |
| Jaeger UI | 16686 | Distributed tracing |
| Prometheus | 9090 | Metrics |
| Grafana | 3000 | Visualization |

## 📊 Architecture Principles

We follow Clean Architecture principles:
- **Domain Layer**: Business logic and entities
- **Application Layer**: Use cases and orchestration
- **Infrastructure Layer**: External services and databases
- **Delivery Layer**: HTTP/gRPC handlers

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the LICENSE file for details.