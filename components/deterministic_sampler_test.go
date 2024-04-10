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

package components

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/trace"
)

func TestDeterministicSamplerSetup(t *testing.T) {
	testCases := []struct {
		name             string
		sampleRate       int
		decision         trace.SamplingDecision
		innersamplerDesc string
	}{
		{
			name:             "negative sample rate",
			sampleRate:       -1,
			decision:         trace.Drop,
			innersamplerDesc: "AlwaysOffSampler",
		},
		{
			name:             "sample rate 0 -- never sample",
			sampleRate:       0,
			decision:         trace.Drop,
			innersamplerDesc: "AlwaysOffSampler",
		},
		{
			name:             "sample rate 1 -- always sample",
			sampleRate:       1,
			decision:         trace.RecordAndSample,
			innersamplerDesc: "AlwaysOnSampler",
		},
		{
			name:             "sample rate 10 -- ratio based",
			sampleRate:       10,
			decision:         trace.RecordAndSample,
			innersamplerDesc: "TraceIDRatioBased{0.1}",
		},
		{
			name:             "sample rate 100 -- ratio based",
			sampleRate:       100,
			decision:         trace.RecordAndSample,
			innersamplerDesc: "TraceIDRatioBased{0.01}",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sampler := NewDeterministicSampler(tc.sampleRate)
			assert.Equal(t, "DeterministicSampler", sampler.Description())
			assert.Equal(t, tc.innersamplerDesc, sampler.innerSampler.Description())

			result := sampler.ShouldSample(trace.SamplingParameters{})
			assert.Equal(t, tc.decision, result.Decision)

			if tc.sampleRate > 0 {
				attr := getAttributeWithKey(result.Attributes, "SampleRate")
				if attr == nil {
					t.Fatalf("SampleRate attribute was not found")
				}
				assert.Equal(t, int64(tc.sampleRate), attr.Value.AsInt64())
			}
		})
	}
}

func getAttributeWithKey(attrs []attribute.KeyValue, key string) *attribute.KeyValue {
	for _, attr := range attrs {
		if attr.Key == attribute.Key("SampleRate") {
			return &attr
		}
	}
	return nil
}
