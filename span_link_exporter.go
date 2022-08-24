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

type spanLinkExporter struct {
	linkUrl string
}

var _ trace.SpanExporter = (*spanLinkExporter)(nil)

func NewSpanLinkExporter(apikey string, serviceName string) (*spanLinkExporter, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.honeycomb.io/1/auth", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Honeycomb-Team", apikey)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var hnyAuthResp honeycombAuthResponse
	err = json.Unmarshal(body, &hnyAuthResp)
	if err != nil {
		return nil, err
	}

	linkUrl := fmt.Sprintf("https://ui.honeycomb.io/%s", hnyAuthResp.Team.Slug)
	if !isClassicApiKey(apikey) {
		linkUrl += fmt.Sprintf("/environments/%s", hnyAuthResp.Environment.Slug)
	}
	linkUrl += fmt.Sprintf("/datasets/%s/trace?trace_id", serviceName)

	return &spanLinkExporter{
		linkUrl: linkUrl,
	}, nil
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

// Export spans is required to implement the Exporter interface.
// It does not actually export spans. Instead, it builds a link to
// honeycomb for the trace that was created, then prints it out!
func (e *spanLinkExporter) ExportSpans(ctx context.Context, spans []trace.ReadOnlySpan) error {
	if len(spans) == 0 {
		return nil
	}

	for _, span := range spans {
		// if a root span (ie no parent span ID)
		if !span.Parent().SpanID().IsValid() {
			fmt.Printf("Trace for %s\nHoneycomb link: %s=%s\n", span.Name(), e.linkUrl, span.SpanContext().TraceID().String())
		}
	}

	return nil
}

// Shutdown is called to stop the exporter, it preforms no action.
func (e *spanLinkExporter) Shutdown(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	return nil
}
