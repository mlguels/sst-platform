package models

import "time"

type Telemetry struct {
 DeviceID string
 Timestamp time.Time
 VoltageV float64
 CurrentA float64
 TemperatureC float64
}

func (t Telemetry) Validate() error
// Makes sure the required fields are present and with the correct constraints

func (t Telemetry) Normalize(now func() time.Time) Telemetry
// If Timestamp is zero, set it to now().
