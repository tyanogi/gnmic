# Technology Stack

## Core Technologies
- **Programming Language**: Go (v1.24.7)
- **Primary Protocols**: gRPC, gNMI (gRPC Network Management Interface)

## Backend & Integration
- **Message Brokers**: Kafka (via IBM/sarama), NATS (including Jetstream)
- **Time Series Databases**: Prometheus, InfluxDB, VictoriaMetrics
- **Key-Value Stores / Locking**: Redis, Consul
- **Infrastructure**: Docker (via Docker SDK), Kubernetes (via lockers)

## CLI & UX
- **Interactive Shell**: go-prompt
- **Configuration**: Support for YAML, TOML, JSON
- **Templating**: Go templates (gtemplate)

## Architecture
- **Structure**: Modular CLI application with a plugin-based architecture for inputs, outputs, and processors.
- **Concurrency**: Extensive use of Go routines for high-performance telemetry collection and processing.
