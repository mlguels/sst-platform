# Telemetry Ingest (/ingest)

# POST /ingest

# Payload (JSON)

- device_id string required
- ts (string, ISO-8601 UTC) optional; if missing, backend sets server time
- voltage_v (number, v) required
- current_a (number, amps) required
- temperature_c (number, °C) — required

# Constraints

- device*id non-empty; max 64 chars (alnum, -, *).
- voltage_v in [0, 1000]; current_a in [0, 1000]; temperature_c in [-40, 150].

# Responses

- 202 Accepted on success
- 400 Bad Request
-

# Headers

- Content-Type: application/json

# Metrics Query (/metrics)

Later

# Command (/command)

Phase 3

# WebSocket (/ws)

Phase 2
