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
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/trace"
)

type DeterministicSampler struct {
	innerSampler        trace.Sampler
	sampleRateAttribute attribute.KeyValue
}

func NewDeterministicSampler(sampleRate int) DeterministicSampler {
	var innerSampler trace.Sampler
	switch {
	case sampleRate <= 0:
		innerSampler = trace.NeverSample()
	case sampleRate == 1:
		innerSampler = trace.AlwaysSample()
	default:
		innerSampler = trace.TraceIDRatioBased(1.0 / float64(sampleRate))
	}
	return DeterministicSampler{
		innerSampler:        innerSampler,
		sampleRateAttribute: attribute.Int("SampleRate", sampleRate),
	}
}

func (ds DeterministicSampler) ShouldSample(parameters trace.SamplingParameters) trace.SamplingResult {
	result := ds.innerSampler.ShouldSample(parameters)
	if result.Decision == trace.RecordAndSample {
		result.Attributes = append(result.Attributes, ds.sampleRateAttribute)
	}
	return result
}

func (ds DeterministicSampler) Description() string {
	return "DeterministicSampler"
}
