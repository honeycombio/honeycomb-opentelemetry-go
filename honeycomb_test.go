package honeycomb

import (
	"fmt"
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
	}
}

func TestSetVendorOptionsWithApikeyAndDataset(t *testing.T) {
	t.Setenv("HONEYCOMB_API_KEY", "atestkey")
	t.Setenv("HONEYCOMB_DATASET", "adataset")

	aConfig := freshConfig()

	for _, setter := range getVendorOptionSetters() {
		setter(aConfig)
	}

	assert.Equal(t, DefaultSpanExporterEndpoint, aConfig.TracesExporterEndpoint)
	assert.Equal(t, DefaultMetricExporterEndpoint, aConfig.MetricsExporterEndpoint)
	assert.Equal(t, "atestkey", aConfig.Headers[honeycombApiKeyHeader])
	assert.Equal(t, "adataset", aConfig.Headers[honeycombDatasetHeader])
}

func TestSetVendorOptionsNoApikeyAndDataset(t *testing.T) {
	t.Setenv("HONEYCOMB_API_KEY", "")
	t.Setenv("HONEYCOMB_DATASET", "")

	aConfig := freshConfig()

	for _, setter := range getVendorOptionSetters() {
		setter(aConfig)
	}

	assert.Equal(t, DefaultSpanExporterEndpoint, aConfig.TracesExporterEndpoint)
	assert.Equal(t, DefaultMetricExporterEndpoint, aConfig.MetricsExporterEndpoint)
	assert.Equal(t, map[string]string{}, aConfig.Headers, "No headers set without config env vars set.")
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
			expectedError: fmt.Errorf("do not include dataset header for non-classic API keys"),
		},
		{
			desc:          "empty API key",
			apikey:        "",
			dataset:       "doesn't matter",
			expectedError: fmt.Errorf("missing x-honeycomb-team header"),
		},
		{
			desc:          "classic API key and no dataset",
			apikey:        "12345678901234567890123456789012",
			dataset:       "",
			expectedError: fmt.Errorf("missing x-honeycomb-dataset header"),
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
