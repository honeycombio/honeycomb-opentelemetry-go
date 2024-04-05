package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/honeycombio/honeycomb-opentelemetry-go"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/trace"
)

// Implement an HTTP Handler function to be instrumented
func httpHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World")
}

// Wrap the HTTP handler function with OTel HTTP instrumentation
func wrapHandler() {
	handler := http.HandlerFunc(httpHandler)
	wrappedHandler := otelhttp.NewHandler(handler, "hello")
	http.Handle("/hello", wrappedHandler)
}

// Running the example:
// export OTEL_SERVICE_NAME=distroless-example
// export OTEL_EXPORTER_OTLP_HEADERS="x-honeycomb-team=XXXX"
// export OTEL_EXPORTER_OTLP_ENDPOINT=https://api.honeycomb.io:443
// export HONEYCOMB_ENABLE_LOCAL_VISUALIZATIONS=1
// go run example.go
func main() {
	var opts []trace.TracerProviderOption
	// Enable local visualizations
	if enableLocalVisualizationsStr := os.Getenv("HONEYCOMB_ENABLE_LOCAL_VISUALIZATIONS"); enableLocalVisualizationsStr != "" {
		var apikey, headers, serviceName string
		if serviceName = os.Getenv("OTEL_SERVICE_NAME"); serviceName == "" {
			serviceName = "unknown_service:go"
		}

		if headers = os.Getenv("OTEL_EXPORTER_OTLP_HEADERS"); headers != "" {
			for _, header := range strings.Split(headers, ",") {
				if strings.HasPrefix(header, "x-honeycomb-team=") {
					if parts := strings.Split(header, "="); len(parts) > 1 {
						apikey = parts[1]
					}
				}
			}
		}
		enabled, _ := strconv.ParseBool(enableLocalVisualizationsStr)
		if enabled {
			exporter, _ := honeycomb.NewSpanLinkExporter(apikey, serviceName)
			opts = append(opts, trace.WithSpanProcessor(trace.NewSimpleSpanProcessor(exporter)))
		}
	}
	// Enable multi-span attributes
	opts = append(opts, trace.WithSpanProcessor(honeycomb.NewBaggageSpanProcessor()))

	if sampleRateStr := os.Getenv("SAMPLE_RATE"); sampleRateStr != "" {
		sampleRate, err := strconv.Atoi(sampleRateStr)
		if err == nil {
			opts = append(opts, trace.WithSampler(honeycomb.NewDeterministicSampler(sampleRate)))
		}
	}

	ctx := context.Background()
	// Configure a new OTLP exporter using environment variables for sending data to Honeycomb over gRPC
	client := otlptracegrpc.NewClient()
	exp, err := otlptrace.New(ctx, client)
	if err != nil {
		log.Fatalf("failed to initialize exporter: %e", err)
	}
	opts = append(opts, trace.WithBatcher(exp))

	// Create a new tracer provider with a batch span processor and the otlp exporter
	tp := trace.NewTracerProvider(opts...)

	// Handle shutdown to ensure all sub processes are closed correctly and telemetry is exported
	defer func() {
		_ = exp.Shutdown(ctx)
		_ = tp.Shutdown(ctx)
	}()

	// Register the global Tracer provider
	otel.SetTracerProvider(tp)

	// Register the W3C trace context and baggage propagators so data is propagated across services/processes
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	// Initialize HTTP handler instrumentation
	wrapHandler()
	port := 3030
	listenAddr := fmt.Sprintf(":%d", port)
	log.Printf("Now listening on:%d ....\n", port)
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}
