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

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/trace"
)

type dynamicAttributeSpanProcessor struct {
	SetAttributes func() []attribute.KeyValue
}

var _ trace.SpanProcessor = (*dynamicAttributeSpanProcessor)(nil)

func NewDynamicAttributeSpanProcessor(setAttributes func() []attribute.KeyValue) trace.SpanProcessor {
	return &dynamicAttributeSpanProcessor{
		SetAttributes: setAttributes,
	}
}

func (processor dynamicAttributeSpanProcessor) OnStart(_ context.Context, span trace.ReadWriteSpan) {
	span.SetAttributes(processor.SetAttributes()...)
}
func (processor dynamicAttributeSpanProcessor) OnEnd(s trace.ReadOnlySpan)       {}
func (processor dynamicAttributeSpanProcessor) Shutdown(context.Context) error   { return nil }
func (processor dynamicAttributeSpanProcessor) ForceFlush(context.Context) error { return nil }
