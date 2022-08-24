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
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"go.opentelemetry.io/otel/sdk/trace"
)

func NewSpanLinkExporter(apikey string, serviceName string) (*Exporter, error) {
	client := &http.Client{}
	req, reqErr := http.NewRequest("GET", "https://api.honeycomb.io/1/auth", nil)
	if reqErr != nil {
		return nil, reqErr
	}

	req.Header.Set("X-Honeycomb-Team", apikey)

	resp, respErr := client.Do(req)
	if respErr != nil {
		return nil, respErr
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var hnyAuthResp honeycombAuthResponse
	json.Unmarshal(b, &hnyAuthResp)

	return &Exporter{
		environmentSlug: hnyAuthResp.Environment.Slug,
		teamSlug:        hnyAuthResp.Team.Slug,
		apiKey:          apikey,
		serviceName:     serviceName,
	}, nil
}

type Exporter struct {
	teamSlug        string
	environmentSlug string
	apiKey          string
	serviceName     string
}

type honeycombAuthResponse struct {
	Environment environment `json:"environment"`
	Team        team        `json:"team"`
}

type environment struct {
	Slug string `json:"slug"`
}

type team struct {
	Slug string `json:"slug"`
}

func getTraceLink(apikey string, teamSlug string, environmentSlug string, serviceName string, traceID string) string {
	if isClassicApiKey(apikey) {
		return fmt.Sprintf("http://ui.honeycomb.io/%s/datasets/%s/trace?trace_id=%s", teamSlug, serviceName, traceID)
	}
	return fmt.Sprintf("http://ui.honeycomb.io/%s/environments/%s/datasets/%s/trace?trace_id=%s", teamSlug, environmentSlug, serviceName, traceID)
}

// Export spans is required to implement the Exporter interface.
// It does not actually export spans. Instead, it builds a link to
// honeycomb for the trace that was created, then prints it out!
func (e *Exporter) ExportSpans(ctx context.Context, spans []trace.ReadOnlySpan) error {
	if len(spans) == 0 {
		return nil
	}

	for i := range spans {
		span := spans[i]

		if !span.Parent().SpanID().IsValid() {
			fmt.Printf("Trace for %s\n", span.Name())
			traceId := span.SpanContext().TraceID().String()
			link := getTraceLink(e.apiKey, e.teamSlug, e.environmentSlug, e.serviceName, traceId)
			fmt.Printf("Honeycomb link: %s\n", link)
		}
	}

	return nil
}

// Shutdown is called to stop the exporter, it preforms no action.
func (e *Exporter) Shutdown(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	return nil
}
