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
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/honeycombio/honeycomb-opentelemetry-go"
	"github.com/honeycombio/otel-launcher-go/launcher"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
)

func main() {
	dsp := honeycomb.NewDynamicAttributeSpanProcessor(func() []attribute.KeyValue {
		return []attribute.KeyValue{
			attribute.String("app.guru_meditation", getGuruMeditation()),
			attribute.Int64("app.unix_time_ms", time.Now().UnixMilli()),
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

	suitcase := baggage.FromContext(ctx)
	packingCube, err := baggage.NewMember("app.luggage", url.QueryEscape("set before bar started"))
	if err != nil {
		fmt.Printf("Invalid baggage member: %s.\n", err)
		os.Exit(1)
	}
	suitcase, err = suitcase.SetMember(packingCube)
	if err != nil {
		fmt.Printf("I couldn't pack that: %s.\n", err)
		os.Exit(1)
	}
	ctx = baggage.ContextWithBaggage(ctx, suitcase)

	ctx, barSpan := tracer.Start(ctx, "bar")
	defer barSpan.End()

	_, bazSpan := tracer.Start(ctx, "baz")
	defer bazSpan.End()

	fmt.Println("OpenTelemetry example")
}

func getGuruMeditation() string {
	bytes := make([]byte, 4)
	if _, err := rand.Read(bytes); err != nil {
		return "48454C50"
	}
	return strings.ToUpper(hex.EncodeToString(bytes))
}
