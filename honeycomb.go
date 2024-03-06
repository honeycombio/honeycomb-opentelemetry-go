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
	"os"
	"regexp"
	"runtime"
	"strconv"

	"github.com/honeycombio/otel-config-go/otelconfig"

	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/trace"
)

const (
	defaultExporterEndpoint          = "api.honeycomb.io:443"
	honeycombApiKeyHeader            = "x-honeycomb-team"
	honeycombDatasetHeader           = "x-honeycomb-dataset"
	honeycombDistroVersionKey        = "honeycomb.distro.version"
	honeycombDistroRuntimeVersionKey = "honeycomb.distro.runtime_version"
	otlpProtoVersionHeader           = "x-otlp-version"
	otlpProtoVersionValue            = "1.0.0"
)

var classicKeyRegex = regexp.MustCompile(`^[a-f0-9]*$`)
var classicIngestKeyRegex = regexp.MustCompile(`^hc[a-z]ic_[a-z0-9]*$`)

func init() {
	otelconfig.SetVendorOptions = getVendorOptionSetters
	otelconfig.ValidateConfig = validateConfig
	otelconfig.DefaultExporterEndpoint = defaultExporterEndpoint
}

// WithHoneycomb() sets the destination for traces and metrics to Honeycomb's API endpoint.
func WithHoneycomb() otelconfig.Option {
	return func(c *otelconfig.Config) {
		c.ResourceAttributes[honeycombDistroVersionKey] = Version
		c.ResourceAttributes[honeycombDistroRuntimeVersionKey] = runtime.Version()
		c.Headers[otlpProtoVersionHeader] = otlpProtoVersionValue
	}
}

// WithApiKey() sets the authorization header appropriately for sending to Honeycomb's API endpoint.
func WithApiKey(apikey string) otelconfig.Option {
	return func(c *otelconfig.Config) {
		c.Headers[honeycombApiKeyHeader] = apikey
	}
}

// WithTracesApiKey() sets the authorization header appropriately for sending traces telemetry to Honeycomb's API endpoint.
func WithTracesApiKey(apikey string) otelconfig.Option {
	return func(c *otelconfig.Config) {
		c.TracesHeaders[honeycombApiKeyHeader] = apikey
	}
}

// WithMetricsApiKey() sets the authorization header appropriately for sending metrics telemetry to Honeycomb's API endpoint.
func WithMetricsApiKey(apikey string) otelconfig.Option {
	return func(c *otelconfig.Config) {
		c.MetricsHeaders[honeycombApiKeyHeader] = apikey
	}
}

// WithDataset() sets the header for routing telemetry to a named dataset at Honeycomb. (For trace data in Classic teams and for metrics only.)
func WithDataset(dataset string) otelconfig.Option {
	return func(c *otelconfig.Config) {
		c.Headers[honeycombDatasetHeader] = dataset
	}
}

// WithTracesDataset() sets the header for routing traces telemetry to a named dataset at Honeycomb.
func WithTracesDataset(dataset string) otelconfig.Option {
	return func(c *otelconfig.Config) {
		c.TracesHeaders[honeycombDatasetHeader] = dataset
	}
}

// WithMetricsDataset() sets the header for routing metrics telemetry to a named dataset at Honeycomb.
func WithMetricsDataset(dataset string) otelconfig.Option {
	return func(c *otelconfig.Config) {
		c.MetricsHeaders[honeycombDatasetHeader] = dataset
	}
}

// WithSampler() sets the sampler used to sample trace spans using a Honeycomb sample rate.
// Sample rate is expressed as 1/X where x is the population size.
func WithSampler(sampleRate int) otelconfig.Option {
	return func(c *otelconfig.Config) {
		c.Sampler = NewDeterministicSampler(sampleRate)
	}
}

// WithDebugSpanExporter() determines whether a debug (stdout) traces exporter should be configured.
func WithDebugSpanExporter() otelconfig.Option {
	spanExporter, _ := stdouttrace.New(stdouttrace.WithPrettyPrint())
	return otelconfig.WithSpanProcessor(trace.NewSimpleSpanProcessor(spanExporter))
}

func getVendorOptionSetters() []otelconfig.Option {
	opts := []otelconfig.Option{
		WithHoneycomb(),
	}

	apikey := ""
	serviceName := "unknown_service:go"

	if endpoint := os.Getenv("HONEYCOMB_API_ENDPOINT"); endpoint != "" {
		opts = append(opts, otelconfig.WithExporterEndpoint(endpoint))
	}
	if endpoint := os.Getenv("HONEYCOMB_TRACES_API_ENDPOINT"); endpoint != "" {
		opts = append(opts, otelconfig.WithTracesExporterEndpoint(endpoint))
	}
	if endpoint := os.Getenv("HONEYCOMB_METRICS_API_ENDPOINT"); endpoint != "" {
		opts = append(opts, otelconfig.WithMetricsExporterEndpoint(endpoint))
	}
	if apikey = os.Getenv("HONEYCOMB_API_KEY"); apikey != "" {
		opts = append(opts, WithApiKey(apikey))
	}

	if apikey := os.Getenv("HONEYCOMB_TRACES_APIKEY"); apikey != "" {
		opts = append(opts, WithTracesApiKey(apikey))
	}
	if apikey := os.Getenv("HONEYCOMB_METRICS_APIKEY"); apikey != "" {
		opts = append(opts, WithMetricsApiKey(apikey))
	}
	if dataset := os.Getenv("HONEYCOMB_DATASET"); dataset != "" {
		opts = append(opts, WithDataset(dataset))
	}

	if dataset := os.Getenv("HONEYCOMB_TRACES_DATASET"); dataset != "" {
		opts = append(opts, WithTracesDataset(dataset))
	}
	if dataset := os.Getenv("HONEYCOMB_METRICS_DATASET"); dataset != "" {
		opts = append(opts, WithMetricsDataset(dataset))
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
			opts = append(opts, otelconfig.WithLogLevel("debug"))
		}
	}

	if serviceName = os.Getenv("OTEL_SERVICE_NAME"); serviceName == "" {
		opts = append(opts, otelconfig.WithServiceName("unknown_service:go"))
	}

	if enableLocalVisualizationsStr := os.Getenv("HONEYCOMB_ENABLE_LOCAL_VISUALIZATIONS"); enableLocalVisualizationsStr != "" {
		enabled, _ := strconv.ParseBool(enableLocalVisualizationsStr)
		if enabled {
			exporter, _ := NewSpanLinkExporter(apikey, serviceName)
			sp := otelconfig.WithSpanProcessor(trace.NewSimpleSpanProcessor(exporter))
			opts = append(opts, sp)
		}
	}

	// default metrics off unless explicity enabled
	metricsEnabled := false
	if enabledStr := os.Getenv("OTEL_METRICS_ENABLED"); enabledStr != "" {
		enabled, _ := strconv.ParseBool(enabledStr)
		if enabled {
			metricsEnabled = true
		}
	}
	opts = append(opts, otelconfig.WithMetricsEnabled(metricsEnabled))
	return opts
}

func validateConfig(c *otelconfig.Config) error {
	apikey := c.Headers[honeycombApiKeyHeader]
	dataset := c.Headers[honeycombDatasetHeader]

	if c.Logger != nil {
		if len(apikey) == 0 {
			c.Logger.Debugf(noApiKeyDetectedMessage)
		} else if isClassicKey(apikey) {
			if dataset == "" {
				c.Logger.Debugf("%s\n%s", classicKeyMissingDatasetMessage, apikey)
			}
		} else {
			if dataset != "" {
				c.Logger.Debugf(dontSetADatasetMessageMessage)
			}
		}
	}

	return nil
}

func isClassicKey(key string) bool {
	if len(key) == 32 {
		return classicKeyRegex.MatchString(key)
	} else if len(key) == 64 {
		return classicIngestKeyRegex.MatchString(key)
	}
	return false
}
