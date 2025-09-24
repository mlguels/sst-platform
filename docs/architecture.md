How to Read This Doc

Top half = quick map.
Skim the Problem, Goals, Components, Interfaces, Data Model to remind me what this system is about

Bottom half = detailed blueprint.
When needed specifics: diagrams, flows, or backend component breakdowns. This is the “how it really works” section.

Use with glossary.
New terms I come across check docs/glossary.md for a definition.

Update as you build.

Add to Interfaces when you finalize a new API.

Adjust Data Model when schema changes.

Mark open questions as resolved.

Purpose: Map so I don’t feel lost when juggling backend, dashboard, and MQTT all at once.

# Architecture — SST Edge Monitoring & Control

## 1) Problem & Goals

- Ingest SST telemetry, visualize it, alert on thresholds, and send control commands safely.

## 2) Non-Goals (for now)

- No real power switching control; device is simulated.
- No multi-tenant auth; use simple API keys.

## 3) System Overview

Device → Backend (/ingest) → TSDB → Analytics → Alerts → UI, and UI → Backend (/command) → Device (via MQTT).

## 4) Components

- **Backend (Go):** ingestion, analytics, control, storage.
- **Dashboard (Next.js):** charts, alerts, controls.
- **TSDB (TimescaleDB or InfluxDB):** telemetry storage.
- **MQTT broker:** southbound messaging.
- **Device simulator:** publishes telemetry, applies commands.

## 5) Interfaces

- **REST:**

  - POST /ingest (telemetry payload)
  - GET /metrics (query window)
  - POST /command (setpoints/modes, idempotency key)

- **WebSocket:** /ws (push latest samples + alerts)
- **MQTT topics:**

  - sst/{deviceId}/telemetry
  - sst/{deviceId}/command

## 6) Data Model

- telemetry(device_id, ts, voltage, current, temperature)
- alerts(id, device_id, ts, type, value, threshold)
- commands(id, device_id, ts, type, payload, status, idem_key)

## 7) Operations

- **Logging:** structured (request_id, device_id)
- **Metrics:** Prometheus counters for requests/errors
- **Reliability:** retries with backoff; idempotent commands

## 8) Open Questions

- TimescaleDB vs InfluxDB?
- WebSockets vs SSE for UI updates?

---

# Detailed Diagrams & Flows

## 1) High-Level System

```mermaid
flowchart LR
  subgraph Field["SST & Field Devices"]
    SST["Solid State Transformer\n(Firmware/RTOS)"]
    Sensors["Voltage/Current/Temp\nSensors"]
    Actuators["Gate Drivers / Relays\n(Controlled by FW)"]
  end

  subgraph Edge["Edge Network"]
    ProtoGW["Protocol Adapter (MQTT/gRPC/REST)\n(optional)"]
  end

  subgraph Backend["Go Backend"]
    API["REST/gRPC API\n(cmd/server)"]
    Ingestion["Ingestion & Parsing\n(internal/ingestion)"]
    Analytics["Rule Engine & Anomaly Detect\n(internal/analytics)"]
    Control["Command Interface\n(internal/control)"]
    Storage["Database Client\n(internal/storage)"]
    Alerts["Alerting/Notifs\n(email/webhook)"]
  end

  subgraph Data["Data Stores"]
    TSDB["Time-Series DB\n(InfluxDB/TimescaleDB)"]
    MetaDB["Metadata/Config\n(Postgres or TSDB ext)"]
  end

  subgraph UI["Dashboard (Next.js)"]
    Viz["Real-time Charts & Logs"]
    Panel["Controls & Thresholds"]
  end

  Sensors -- telemetry --> SST
  SST -- telemetry --> ProtoGW
  ProtoGW -- HTTP/gRPC --> API
  API --> Ingestion --> Storage
  Storage <--> TSDB
  Analytics --> Alerts
  UI <-- WebSocket/HTTP --> API
  UI -.-> Panel -.-> Control
  Control -- command --> ProtoGW -- control --> SST
```

## 2) Backend Components (Go)

```mermaid
flowchart TB
  subgraph cmd/server
    Main["main.go\nwire up routes, config, logging"]
  end

  subgraph pkg/api
    Routes["HTTP/gRPC routes\n/ingest /metrics /command\n/ws (live)"]
    Schemas["DTOs & validation"]
  end

  subgraph internal/ingestion
    Parser["Parsers (MQTT/JSON/protobuf)"]
    Normalizer["Units, bounds, timestamps"]
  end

  subgraph internal/analytics
    Rules["Threshold checks"]
    Anomaly["Rolling stats, EWMA,\nrate-of-change guards"]
  end

  subgraph internal/control
    CmdBus["Command dispatcher\n(idempotency, audit)"]
    Adapters["MQTT/gRPC clients to field"]
  end

  subgraph internal/storage
    Repo["Repository interfaces"]
    TSClient["TSDB client"]
  end

  subgraph pkg/utils
    Log["logger"]
    Cfg["config loader"]
    Err["error types"]
  end

  Main --> Routes
  Routes --> Parser --> Normalizer --> Repo --> TSClient
  Repo --> TSClient
  Normalizer --> Rules --> Anomaly
  Routes --> CmdBus --> Adapters
  Rules --> Alerts[internal/alerts]
  Anomaly --> Alerts
  Routes <---> WS["WebSocket hub\npush updates to UI"]
  Routes --> Schemas
  Main --> Log & Cfg & Err
```

## 3) Telemetry Ingest — Sequence

```mermaid
sequenceDiagram
  autonumber
  participant FW as SST Firmware
  participant GW as Protocol Adapter (MQTT/gRPC)
  participant API as Backend API (/ingest)
  participant ING as Ingestion
  participant ANA as Analytics
  participant DB as TSDB
  participant WS as WS Hub
  participant UI as Dashboard

  FW->>GW: Telemetry (V/I/Temp, ts)
  GW->>API: POST /ingest (JSON/Protobuf)
  API->>ING: Validate, parse, normalize (units, ts)
  ING->>DB: Write time-series points
  ING->>ANA: Publish sample to rule engine
  ANA-->>API: Events (ok/warn/alarm)
  API-->>WS: Push latest sample + events
  WS-->>UI: Live charts + alerts
```

## 4) Control Command — Sequence

```mermaid
sequenceDiagram
  autonumber
  participant UI as Dashboard
  participant API as Backend API (/command)
  participant CTRL as Control Dispatcher
  participant GW as Protocol Adapter (MQTT/gRPC)
  participant FW as SST Firmware
  participant DB as TSDB

  UI->>API: POST /command {setpoint: ...}
  API->>CTRL: Validate, authz, idempotency key
  CTRL->>GW: Publish command (gRPC/MQTT topic)
  GW->>FW: Apply setpoint/mode
  FW-->>GW: Ack + new telemetry
  GW-->>API: POST /ingest (updated state)
  API->>DB: Store command audit + result
  API-->>UI: 202 Accepted (then status via WS/poll)
```
