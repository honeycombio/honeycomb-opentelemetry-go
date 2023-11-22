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
	"context"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/sdk/trace"
)

var _ trace.SpanExporter = &testExporter{}

type testExporter struct {
	spans []trace.ReadOnlySpan
}

func (e *testExporter) Start(ctx context.Context) error    { return nil }
func (e *testExporter) Shutdown(ctx context.Context) error { return nil }

func (e *testExporter) ExportSpans(ctx context.Context, ss []trace.ReadOnlySpan) error {
	e.spans = append(e.spans, ss...)
	return nil
}

func NewTestExporter() *testExporter {
	return &testExporter{}
}

func TestBaggageSpanProcessorAppendsBaggageAttributes(t *testing.T) {
	// create ctx with some baggage
	ctx := context.Background()
	suitcase := baggage.FromContext(ctx)
	packingCube, _ := baggage.NewMember("baggage.test", url.PathEscape("baggage value"))
	suitcase, _ = suitcase.SetMember(packingCube)
	ctx = baggage.ContextWithBaggage(ctx, suitcase)

	// create trace provider with baggage processor and test exporter
	exporter := NewTestExporter()
	tp := trace.NewTracerProvider(
		trace.WithSpanProcessor(NewBaggageSpanProcessor()),
		trace.WithSpanProcessor(trace.NewSimpleSpanProcessor(exporter)),
	)

	// create tracer and start/end span
	tracer := tp.Tracer("test")
	_, span := tracer.Start(ctx, "test")
	span.End()

	assert.Equal(t, 1, len(exporter.spans))
	assert.Equal(t, 1, len(exporter.spans[0].Attributes()))

	for _, attr := range exporter.spans[0].Attributes() {
		assert.Equal(t, attribute.Key("baggage.test"), attr.Key)
		assert.Equal(t, "baggage value", attr.Value.AsString())
	}
}
