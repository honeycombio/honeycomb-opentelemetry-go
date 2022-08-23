// Copyright Honeycomb Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package honeycomb

import (
	"fmt"
	"os"
	"runtime"
	"strconv"

	"github.com/honeycombio/opentelemetry-go-contrib/launcher"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/trace"
)

const (
	DefaultSpanExporterEndpoint      = "api.honeycomb.io:443"
	DefaultMetricExporterEndpoint    = "api.honeycomb.io:443"
	honeycombApiKeyHeader            = "x-honeycomb-team"
	honeycombDatasetHeader           = "x-honeycomb-dataset"
	honeycombDistroVersionKey        = "honeycomb.distro.version"
	honeycombDistroRuntimeVersionKey = "honeycomb.distro.runtime_version"
	otlpProtoVersionHeader           = "x-otlp-version"
	otlpProtoVersionValue            = "0.18.0"
)

func init() {
	launcher.SetVendorOptions = getVendorOptionSetters
	launcher.ValidateConfig = validateConfig
}

// WithHoneycomb() sets the destination for traces and metrics to Honeycomb's API endpoint.
func WithHoneycomb() launcher.Option {
	return func(c *launcher.Config) {
		c.ResourceAttributes[honeycombDistroVersionKey] = Version
		c.ResourceAttributes[honeycombDistroRuntimeVersionKey] = runtime.Version()
		c.Headers[otlpProtoVersionHeader] = otlpProtoVersionValue
		c.ExporterEndpoint = DefaultSpanExporterEndpoint
	}
}

// WithApiKey() sets the authorization header appropriately for sending to Honeycomb's API endpoint.
func WithApiKey(apikey string) launcher.Option {
	return func(c *launcher.Config) {
		c.Headers[honeycombApiKeyHeader] = apikey
	}
}

// WithDataset() sets the header for routing telemetry to a named dataset at Honeycomb. (For trace data in Classic teams and for metrics only.)
func WithDataset(dataset string) launcher.Option {
	return func(c *launcher.Config) {
		c.Headers[honeycombDatasetHeader] = dataset
	}
}

// WithSampler() sets the sampler used to sample trace spans using a Honeycomb sample rate.
// Sample rate is expressed as 1/X where x is the population size.
func WithSampler(sampleRate int) launcher.Option {
	return func(c *launcher.Config) {
		c.Sampler = NewDeterministicSampler(sampleRate)
	}
}

// WithDebugSpanExporter() determines whether a debug (stdout) traces exporter should be configured.
func WithDebugSpanExporter() launcher.Option {
	spanExporter, _ := stdouttrace.New(stdouttrace.WithPrettyPrint())
	return launcher.WithSpanProcessor(trace.NewSimpleSpanProcessor(spanExporter))
}

func getVendorOptionSetters() []launcher.Option {
	opts := []launcher.Option{
		WithHoneycomb(),
	}
	if endpoint := os.Getenv("HONEYCOMB_API_ENDPOINT"); endpoint != "" {
		opts = append(opts, launcher.WithExporterEndpoint(endpoint))
	}
	if endpoint := os.Getenv("HONEYCOMB_TRACES_API_ENDPOINT"); endpoint != "" {
		opts = append(opts, launcher.WithSpanExporterEndpoint(endpoint))
	}
	if endpoint := os.Getenv("HONEYCOMB_METRICS_API_ENDPOINT"); endpoint != "" {
		opts = append(opts, launcher.WithMetricExporterEndpoint(endpoint))
	}
	if apikey := os.Getenv("HONEYCOMB_API_KEY"); apikey != "" {
		opts = append(opts, WithApiKey(apikey))
	}
	if dataset := os.Getenv("HONEYCOMB_DATASET"); dataset != "" {
		opts = append(opts, WithDataset(dataset))
	}
	if sampleRateStr := os.Getenv("SAMPLE_RATE"); sampleRateStr != "" {
		sampleRate, err := strconv.Atoi(sampleRateStr)
		if err == nil {
			opts = append(opts, WithSampler(sampleRate))
		}
	}
	if enabledStr := os.Getenv("DEBUG"); enabledStr != "" {
		enabled, _ := strconv.ParseBool(enabledStr)
		if enabled {
			opts = append(opts, WithDebugSpanExporter())
			opts = append(opts, launcher.WithLogLevel("debug"))
		}
	}
	if serviceName := os.Getenv("OTEL_SERVICE_NAME"); serviceName == "" {
		opts = append(opts, launcher.WithServiceName("unknown_service:go"))
	}
	// default metrics off unless explicity enabled
	metricsEnabled := false
	if enabledStr := os.Getenv("OTEL_METRICS_ENABLED"); enabledStr != "" {
		enabled, _ := strconv.ParseBool(enabledStr)
		if enabled {
			metricsEnabled = true
		}
	}
	opts = append(opts, launcher.WithMetricsEnabled(metricsEnabled))
	return opts
}

func validateConfig(c *launcher.Config) error {
	apikey := c.Headers[honeycombApiKeyHeader]
	dataset := c.Headers[honeycombDatasetHeader]

	switch len(apikey) {
	case 0:
		return fmt.Errorf(noApiKeyDetectedMessage)
	case 32: // classic
		if dataset == "" {
			return fmt.Errorf(classicKeyMissingDatasetMessage, apikey)
		}
	default:
		if dataset != "" {
			return fmt.Errorf(dontSetADatasetMessageMessage)
		}
	}
	return nil
}
