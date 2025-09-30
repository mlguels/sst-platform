-- telemetry table

telemetry(device_id TEXT NOT NULL, ts TIMESTAMPTZ NOT NULL, voltage_v DOUBLE PRECISION NOT NULL, current_a DOUBLE PRECISION NOT NULL, temperature_c DOUBLE PRECISION NOT NULL, PRIMARY KEY (device_id, ts))

