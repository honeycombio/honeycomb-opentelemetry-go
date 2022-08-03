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

package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/honeycombio/opentelemetry-go-contrib/launcher"
	honeycomb "github.com/honeycombio/otel/honeycomb"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func main() {
	dsp := honeycomb.NewDynamicAttributeSpanProcessor(func() []attribute.KeyValue {
		return []attribute.KeyValue{
			attribute.String("foo", "bar"),
			attribute.Int64("unix", time.Now().UnixMilli()),
		}
	})
	bsp := honeycomb.NewBaggageSpanProcessor()

	shutdown, err := launcher.ConfigureOpenTelemetry(
		launcher.WithSpanProcessor(dsp, bsp),
	)
	defer shutdown()

	if err != nil {
		fmt.Printf("Nope. That configuration doesn't seem right: %s.\n", err)
		os.Exit(1)
	}

	tracer := otel.Tracer("honeycomb-otel-go-distro-example")

	ctx := context.Background()

	ctx, fooSpan := tracer.Start(ctx, "foo")
	defer fooSpan.End()

	ctx, barSpan := tracer.Start(ctx, "bar")
	defer barSpan.End()

	_, bazSpan := tracer.Start(ctx, "baz")
	defer bazSpan.End()

	fmt.Println("OpenTelemetry example")
}
