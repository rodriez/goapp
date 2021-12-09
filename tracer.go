package goapp

import (
	"os"

	"github.com/rodriez/gotracing"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var Tracer trace.Tracer
var provider gotracing.TracingProvider

func InitTracing() {
	providerName := os.Getenv("TRACE_PROVIDER")
	traceId := os.Getenv("TRACE_ID")

	Tracer = otel.Tracer(traceId)

	provider = gotracing.Build(providerName)
	provider.Setup(traceId)
}

func CloseTracing() {
	provider.Close()
}
