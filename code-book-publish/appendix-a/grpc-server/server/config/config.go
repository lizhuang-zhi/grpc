package config

import (
	"github.com/practicalgo/code/appendix-a/grpc-server/server/telemetry"
	"github.com/rs/zerolog"
)

type AppConfig struct {
	Logger  zerolog.Logger
	Metrics telemetry.MetricReporter
}
