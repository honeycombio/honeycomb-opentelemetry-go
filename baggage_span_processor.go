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
	"go.opentelemetry.io/contrib/processors/baggage/baggagetrace"
	"go.opentelemetry.io/otel/sdk/trace"
)

// Returns a new baggageSpanProcessor.
//
// The Baggage span processor duplicates onto a span the attributes found
// in Baggage in the parent context at the moment the span is started.
//
// Deprecated: Use go.opentelemetry.io/contrib/processors/baggage/baggagetrace.New instead
func NewBaggageSpanProcessor() trace.SpanProcessor {
	return baggagetrace.New()
}
