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
	"runtime"
	"testing"

	"github.com/honeycombio/otel-launcher-go/launcher"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
)

func freshConfig() *launcher.Config {
	return &launcher.Config{
		TracesExporterEndpoint:          "",
		TracesExporterEndpointInsecure:  false,
		TracesEnabled:                   false,
		ServiceName:                     "",
		ServiceVersion:                  "",
		Headers:                         map[string]string{},
		TracesHeaders:                   map[string]string{},
		MetricsHeaders:                  map[string]string{},
		MetricsExporterEndpoint:         "",
		MetricsExporterEndpointInsecure: false,
		MetricsEnabled:                  false,
		MetricsReportingPeriod:          "",
		LogLevel:                        "",
		Propagators:                     []string{},
		ResourceAttributes:              map[string]string{},
		ResourceAttributesFromEnv:       "",
		SpanProcessors:                  []trace.SpanProcessor{},
		Resource:                        &resource.Resource{},
		Logger:                          nil,
		ShutdownFunctions:               []func(c *launcher.Config) error{},
		Sampler:                         trace.AlwaysSample(),
	}
}

func TestSetVendorOptions(t *testing.T) {
	testCases := []struct {
		desc            string
		apikey          string
		dataset         string
		expectedHeaders map[string]string
	}{
		{
			desc:    "with API key and dataset",
			apikey:  "atestkey",
			dataset: "adataset",
			expectedHeaders: map[string]string{
				honeycombApiKeyHeader:  "atestkey",
				honeycombDatasetHeader: "adataset",
				otlpProtoVersionHeader: otlpProtoVersionValue,
			},
		},
		{
			desc:    "no API key or dataset",
			apikey:  "",
			dataset: "",
			expectedHeaders: map[string]string{
				otlpProtoVersionHeader: otlpProtoVersionValue,
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			t.Setenv("HONEYCOMB_API_KEY", tC.apikey)
			t.Setenv("HONEYCOMB_DATASET", tC.dataset)

			aConfig := freshConfig()

			for _, setter := range getVendorOptionSetters() {
				setter(aConfig)
			}

			assert.Equal(t, DefaultSpanExporterEndpoint, aConfig.ExporterEndpoint,
				"Trace & metric data should be configured to target the Honeycomb API endpoint.",
			)
			assert.Equal(t, tC.expectedHeaders, aConfig.Headers)
		})
	}
}

func TestValidateConfig(t *testing.T) {
	testCases := []struct {
		desc                 string
		apikey               string
		dataset              string
		expectedLoggerFormat string
		expectedLoggerValues []interface{}
	}{
		{
			desc:                 "modern API key and no dataset",
			apikey:               "123456789012345678901",
			dataset:              "",
			expectedLoggerFormat: "",
			expectedLoggerValues: nil,
		},
		{
			desc:                 "classic API key and a dataset",
			apikey:               "12345678901234567890123456789012",
			dataset:              "is-set-horrah",
			expectedLoggerFormat: "",
			expectedLoggerValues: nil,
		},
		{
			desc:                 "modern API key and a dataset",
			apikey:               "123456789012345678901",
			dataset:              "no thank you",
			expectedLoggerFormat: dontSetADatasetMessageMessage,
			expectedLoggerValues: nil,
		},
		{
			desc:                 "empty API key",
			apikey:               "",
			dataset:              "doesn't matter",
			expectedLoggerFormat: noApiKeyDetectedMessage,
			expectedLoggerValues: nil,
		},
		{
			desc:                 "classic API key and no dataset",
			apikey:               "12345678901234567890123456789012",
			dataset:              "",
			expectedLoggerFormat: "%s\n%s",
			expectedLoggerValues: []interface{}{
				classicKeyMissingDatasetMessage,
				"12345678901234567890123456789012",
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			aConfig := freshConfig()
			aConfig.Headers[honeycombApiKeyHeader] = tC.apikey
			aConfig.Headers[honeycombDatasetHeader] = tC.dataset

			logger := &captureLogger{}

			aConfig.Logger = logger

			err := validateConfig(aConfig)

			// Check the output
			assert.Equal(t, err, nil)
			assert.Equal(t, tC.expectedLoggerFormat, logger.Format)
			assert.Equal(t, tC.expectedLoggerValues, logger.Values)
		})
	}
}

type captureLogger struct {
	Format string
	Values []interface{}
}

func (l *captureLogger) Fatalf(format string, v ...interface{}) {
	l.Format = format
	l.Values = v
}
func (l *captureLogger) Debugf(format string, v ...interface{}) {
	l.Format = format
	l.Values = v
}

func TestHoneycombResourceAttributesAreSet(t *testing.T) {
	config := freshConfig()
	for _, setter := range getVendorOptionSetters() {
		setter(config)
	}

	assert.Equal(t, Version, config.ResourceAttributes["honeycomb.distro.version"])
	assert.Equal(t, runtime.Version(), config.ResourceAttributes["honeycomb.distro.runtime_version"])
}

func TestConfigureDeterministicSampler(t *testing.T) {
	// no env var - should use default sampler
	config := freshConfig()
	for _, setter := range getVendorOptionSetters() {
		setter(config)
	}
	assert.Equal(t, "AlwaysOnSampler", config.Sampler.Description())

	// set env var - should have deterministic sampler
	t.Setenv("SAMPLE_RATE", "1")
	config = freshConfig()
	for _, setter := range getVendorOptionSetters() {
		setter(config)
	}
	assert.Equal(t, "DeterministicSampler", config.Sampler.Description())
}

func TestSettingExportersAddsDebugExporter(t *testing.T) {
	config := freshConfig()
	t.Setenv("DEBUG", "true")

	for _, setter := range getVendorOptionSetters() {
		setter(config)
	}

	// it's really tought to determine if a simple span processor (private type)
	// wrapping a stdouttrace span exporter has been configured
	// Let's check we have at least configured a span processor for now
	assert.Equal(t, 1, len(config.SpanProcessors))
}

func TestSettingExportersAddsLocalVizExporter(t *testing.T) {
	config := freshConfig()
	t.Setenv("HONEYCOMB_ENABLE_LOCAL_VISUALIZATIONS", "true")

	for _, setter := range getVendorOptionSetters() {
		setter(config)
	}

	assert.Equal(t, 1, len(config.SpanProcessors))
}

func TestServiceNameDefaultsToUnknownServiceWhenNotSet(t *testing.T) {
	config := freshConfig()

	// If you are running stuff locally and set this, it will mess up the test.
	// So, we explicit set the env var to be empty.
	t.Setenv("OTEL_SERVICE_NAME", "")
	for _, setter := range getVendorOptionSetters() {
		setter(config)
	}
	assert.Equal(t, "unknown_service:go", config.ServiceName)
}

func TestSettingDebugAlsoSetsLogLevelToDebug(t *testing.T) {
	t.Setenv("DEBUG", "true")
	launcher.ValidateConfig = func(c *launcher.Config) error {
		assert.Equal(t, c.LogLevel, "debug")
		return nil
	}
	_, err := launcher.ConfigureOpenTelemetry()
	assert.Nil(t, err)
}

func TestCanSetEndpointsUsingHoneycombEnvVars(t *testing.T) {
	t.Setenv("HONEYCOMB_API_ENDPOINT", "generic-endpoint")
	t.Setenv("HONEYCOMB_TRACES_API_ENDPOINT", "traces-endpoint")
	t.Setenv("HONEYCOMB_METRICS_API_ENDPOINT", "metrics-endpoint")

	launcher.ValidateConfig = func(c *launcher.Config) error {
		assert.Equal(t, "generic-endpoint", c.ExporterEndpoint)
		assert.Equal(t, "traces-endpoint", c.TracesExporterEndpoint)
		assert.Equal(t, "metrics-endpoint", c.MetricsExporterEndpoint)
		return nil
	}
	_, err := launcher.ConfigureOpenTelemetry()
	assert.Nil(t, err)
}

func TestCanSetTracesAndMetricsSpecificHeaders(t *testing.T) {
	t.Setenv("HONEYCOMB_TRACES_APIKEY", "traces-apikey")
	t.Setenv("HONEYCOMB_TRACES_DATASET", "traces-dataset")
	t.Setenv("HONEYCOMB_METRICS_APIKEY", "metrics-apikey")
	t.Setenv("HONEYCOMB_METRICS_DATASET", "metrics-dataset")

	launcher.ValidateConfig = func(c *launcher.Config) error {
		assert.Equal(t, "traces-apikey", c.TracesHeaders[honeycombApiKeyHeader])
		assert.Equal(t, "traces-dataset", c.TracesHeaders[honeycombDatasetHeader])
		assert.Equal(t, "metrics-apikey", c.MetricsHeaders[honeycombApiKeyHeader])
		assert.Equal(t, "metrics-dataset", c.MetricsHeaders[honeycombDatasetHeader])
		return nil
	}
	_, err := launcher.ConfigureOpenTelemetry()
	assert.Nil(t, err)
}

func TestMetricsAreDisabledByDefault(t *testing.T) {
	// disabled by default
	launcher.ValidateConfig = func(c *launcher.Config) error {
		assert.False(t, c.MetricsEnabled)
		return nil
	}
	_, err := launcher.ConfigureOpenTelemetry()
	assert.Nil(t, err)

	// can be enabled
	t.Setenv("OTEL_METRICS_ENABLED", "true")
	launcher.ValidateConfig = func(c *launcher.Config) error {
		assert.True(t, c.MetricsEnabled)
		return nil
	}
	_, err = launcher.ConfigureOpenTelemetry()
	assert.Nil(t, err)
}
