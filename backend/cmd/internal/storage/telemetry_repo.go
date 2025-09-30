package storage

import "github.com/mlguels/sst-platform/backend/pkg/models"


func InsertSamples(ctx, samples []models.Telemetry) error
func Latest(ctx, deviceID string, limit int) ([]models.Telemetry, error)