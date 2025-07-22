# Architecture Decision Records (ADR)

## ADR-001: Microservices Architecture
**Date**: [Date]
**Status**: Accepted

### Context
Need to build scalable URL shortener with separate concerns.

### Decision
Use 3 microservices: URL, User, Analytics

### Consequences
- ✅ Independent deployment
- ✅ Technology flexibility
- ❌ Increased complexity
- ❌ Network latency

---

## ADR-002: Clean Architecture
**Date**: [Date]
**Status**: Accepted

### Context
Need maintainable and testable code structure.

### Decision
Use Clean Architecture with layers:
- delivery (transport)
- application (use cases)
- domain (business logic)
- infrastructure (external services)

### Consequences
- ✅ Testable business logic
- ✅ Framework independence
- ❌ More boilerplate code
- ❌ Learning curve

---

## ADR-003: Database per Service
**Date**: [Date]
**Status**: Accepted

### Context
Microservices need data isolation.

### Decision
- URL Service: PostgreSQL database
- User Service: PostgreSQL database
- Analytics Service: MongoDB
- Shared Redis for caching

### Consequences
- ✅ Service independence
- ✅ Appropriate DB for use case
- ❌ No ACID across services
- ❌ Data duplication

---

## ADR-004: Message Broker Strategy
**Date**: [Date]
**Status**: Accepted

### Context
Need async communication between services.

### Decision
Phased approach:
1. NATS - lightweight pub/sub
2. RabbitMQ - reliable queues
3. Kafka - event streaming

### Consequences
- ✅ Learn different patterns
- ✅ Right tool for right job
- ❌ Multiple technologies
- ❌ Operational complexity

---

## ADR-005: API Strategy
**Date**: [Date]
**Status**: Proposed

### Context
Need both public and internal APIs.

### Decision
- Public API: REST (HTTP/JSON)
- Internal: gRPC
- Real-time: WebSocket

### Consequences
- ✅ Best protocol for use case
- ✅ Type safety with gRPC
- ❌ Multiple protocols to maintain

---

## Template for new ADRs

## ADR-XXX: Title
**Date**: [Date]
**Status**: Proposed/Accepted/Deprecated

### Context
What is the issue we're seeing that motivates this decision?

### Decision
What is the change that we're proposing/doing?

### Consequences
What becomes easier or more difficult because of this change?