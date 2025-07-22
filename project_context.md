# Project Context for Claude AI

## 🎯 Project Overview
We are building a production-ready URL shortener using Go microservices architecture.

## 👥 Team
- Developer: [Your Name]
- AI Assistant: Claude (Code + Chat)
- Role: Learning Go and microservices patterns

## 🏗️ Architecture Decisions

### Microservices (3 services)
1. **URL Service**: Core shortening logic
2. **User Service**: Auth, rate limiting
3. **Analytics Service**: Click tracking

### Technology Stack
- **Language**: Go 1.22+ (using generics where appropriate)
- **Architecture**: Clean Architecture with modifications
- **Databases**:
    - PostgreSQL (URL + User data)
    - Redis (Caching + Rate limiting)
    - MongoDB (Analytics)
- **Message Brokers** (Phased):
    - Week 1-2: NATS
    - Week 3-4: RabbitMQ
    - Week 5-6: Kafka

### Clean Architecture Layers
```
delivery/      → HTTP/gRPC handlers (transport)
application/   → Use cases, DTOs
domain/        → Entities, business rules
infrastructure/→ DB, external services
```

## 📋 Implementation Plan

### Phase 1: Foundation ✅
- [x] Project structure
- [x] Docker environment
- [x] Basic documentation
- [ ] Proto definitions
- [ ] Shared packages

### Phase 2: URL Service
- [ ] Domain models
- [ ] PostgreSQL repository
- [ ] Redis caching
- [ ] HTTP handlers
- [ ] gRPC server
- [ ] Unit tests
- [ ] Integration tests

### Phase 3: User Service
- [ ] JWT authentication
- [ ] OAuth2 integration
- [ ] Rate limiting
- [ ] API keys
- [ ] RBAC

### Phase 4: Analytics Service
- [ ] Event collection
- [ ] MongoDB integration
- [ ] Real-time WebSocket
- [ ] Reporting

### Phase 5: Integration
- [ ] NATS messaging
- [ ] Service discovery
- [ ] Circuit breakers
- [ ] Distributed tracing

### Phase 6: Advanced
- [ ] RabbitMQ queues
- [ ] Kafka streaming
- [ ] Temporal workflows
- [ ] Performance optimization

## 🔧 Development Guidelines

### Code Style
- Clean Architecture principles
- Domain-Driven Design where appropriate
- Test coverage > 80%
- Meaningful commit messages

### Naming Conventions
- Packages: lowercase, no underscores
- Interfaces: Suffix with -er (Reader, Writer)
- Structs: PascalCase
- Constants: UPPER_SNAKE_CASE

### Testing Strategy
- Unit tests for business logic
- Integration tests with testcontainers
- E2E tests for critical paths
- Benchmarks for performance-critical code
- Use testify/assert and testify/require for assertions

## 📝 Current Status
**Date**: [Update this]
**Current Phase**: Phase 1 - Foundation
**Blockers**: Claude Code timeout issues (using manual file creation)
**Next Steps**: Create proto definitions and shared packages

## 🤔 Open Questions
1. Should we use schema-per-service in PostgreSQL?
2. How to handle distributed transactions?
3. Best approach for service discovery?

## 📚 Learning Goals
- Master Go patterns and idioms
- Understand microservices communication
- Learn distributed systems patterns
- Production-ready code practices