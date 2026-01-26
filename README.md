# Distributed Configuration System

## Overview
This project is a **Distributed Configuration System** consisting of three main services: **Controller**, **Agent**, and **Worker**.  
The Controller manages configuration as the single source of truth, Agents fetch and cache configuration, and Workers consume configuration for business logic.

The system is designed to be **scalable**, **fault-tolerant**, and **simple**, with the Controller kept **standalone**.

---

## Technologies Used
- **Go** – main programming language
- **SQLite** – configuration storage (Controller)
- **Redis Pub/Sub** – configuration update notification (optional trigger)
- **net/http** – HTTP server and client
- **Docker & Docker Compose** – service orchestration

---

## Key Design Concepts
- **Clean Architecture**
  - Separation of handler, usecase, repository, and domain layers
  - Business logic independent from framework and infrastructure

- **Polling with ETag**
  - Agent polls Controller using HTTP conditional requests
  - `304 Not Modified` returned when configuration is unchanged

- **Exponential Backoff**
  - Applied on Agent polling retry when errors occur
  - Prevents overwhelming the Controller during failures

- **Redis Pub/Sub (Optional Enhancement)**
  - Used only as a trigger signal
  - Actual configuration is always pulled via HTTP

---

## How to Run Services (Local)

### 1. Controller
```bash
go run cmd/controller/main.go
```

### 2. Agent
```bash
go run cmd/agent/main.go
```

### 3. Worker
```bash
go run cmd/worker/main.go
```

---

## Build & Run Using Docker Compose

### Build and Start Controlller Services
```bash
docker-compose up --build
```

### Build and Start Agent & Worker Services
```bash
docker-compose up --build
```

### Stop Services
```bash
docker-compose down
```

Each service is built and run independently using its own Dockerfile.

---

## Summary
This project demonstrates a practical and scalable approach to distributed configuration management using Go, with clean architecture, efficient polling, and optional event-driven optimization.
