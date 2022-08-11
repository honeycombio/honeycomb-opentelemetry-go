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
