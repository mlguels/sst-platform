High-Level System

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
Storage["Time-Series DB Client\n(internal/storage)"]
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

Backend Components
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

Telemetry Ingest -- Sequence

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

Control Command -- Sequence

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
