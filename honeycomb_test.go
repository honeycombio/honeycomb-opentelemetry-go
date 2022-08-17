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
	"runtime"
	"testing"

	"github.com/honeycombio/opentelemetry-go-contrib/launcher"
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
		HeadersFromEnv:                  "",
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
			},
		},
		{
			desc:            "no API key or dataset",
			apikey:          "",
			dataset:         "",
			expectedHeaders: map[string]string{},
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

			assert.Equal(t, DefaultSpanExporterEndpoint, aConfig.TracesExporterEndpoint,
				"Trace data should be configured to target the Honeycomb API endpoint.",
			)
			assert.Equal(t, DefaultMetricExporterEndpoint, aConfig.MetricsExporterEndpoint,
				"Metric data should be configured to target the Honeycomb API endpoint.",
			)
			assert.Equal(t, tC.expectedHeaders, aConfig.Headers)
		})
	}
}

func TestValidateConfig(t *testing.T) {
	testCases := []struct {
		desc          string
		apikey        string
		dataset       string
		expectedError error
	}{
		{
			desc:          "modern API key and no dataset",
			apikey:        "123456789012345678901",
			dataset:       "",
			expectedError: nil,
		},
		{
			desc:          "classic API key and a dataset",
			apikey:        "12345678901234567890123456789012",
			dataset:       "is-set-horrah",
			expectedError: nil,
		},
		{
			desc:          "modern API key and a dataset",
			apikey:        "123456789012345678901",
			dataset:       "no thank you",
			expectedError: fmt.Errorf(dontSetADatasetMessageMessage),
		},
		{
			desc:          "empty API key",
			apikey:        "",
			dataset:       "doesn't matter",
			expectedError: fmt.Errorf(noApiKeyDetectedMessage),
		},
		{
			desc:          "classic API key and no dataset",
			apikey:        "12345678901234567890123456789012",
			dataset:       "",
			expectedError: fmt.Errorf(classicKeyMissingDatasetMessage, "12345678901234567890123456789012"),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			aConfig := freshConfig()
			aConfig.Headers[honeycombApiKeyHeader] = tC.apikey
			aConfig.Headers[honeycombDatasetHeader] = tC.dataset

			err := validateConfig(aConfig)
			assert.Equal(t, tC.expectedError, err)
		})
	}
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

func TestSettingExporterDebugEnabledAddsDebugExporter(t *testing.T) {
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

func TestServiceNameDefaultsToUnknownServiceWhenNotSet(t *testing.T) {
	config := freshConfig()
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
