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

type localVisualizationsSpanProcessor struct {
	teamSlug        string
	environmentSlug string
	apiKey          string
	serviceName     string
}

type honeycombAuthResponse struct {
	ApiKeyAccess string      `json:"api_key_access"`
	Environment  environment `json:"environment"`
	Team         team        `json:"team"`
}

type apiKeyAcces struct {
	Events         bool `json:"events"`
	Markers        bool `json:"markers"`
	Triggers       bool `json:"triggers"`
	Boards         bool `json:"boards"`
	Queries        bool `json:"queries"`
	Columns        bool `json:"columns"`
	CreateDatasets bool `json:"createDatasets"`
}

type environment struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type team struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

var _ trace.SpanProcessor = (*baggageSpanProcessor)(nil)

// Returns a new localVisualizationsSpanProcessor.
func NewLocalVisualizationsSpanProcessor(apikey string, serviceName string, endpoint string) trace.SpanProcessor {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://api.honeycomb.io/1/auth", nil)
	req.Header.Set("X-Honeycomb-Team", apikey)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Print(err)
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var hnyAuthResp honeycombAuthResponse
	json.Unmarshal(b, &hnyAuthResp)

	return &localVisualizationsSpanProcessor{
		environmentSlug: hnyAuthResp.Environment.Slug,
		teamSlug:        hnyAuthResp.Team.Slug,
		apiKey:          apikey,
		serviceName:     serviceName,
	}
}

func (processor localVisualizationsSpanProcessor) OnStart(ctx context.Context, span trace.ReadWriteSpan) {
}

func (processor localVisualizationsSpanProcessor) OnEnd(s trace.ReadOnlySpan) {
	if !s.Parent().IsValid() {
		fmt.Printf("Trace for %s\n", s.Name())
		traceId := s.SpanContext().TraceID().String()
		fmt.Printf("Honeycomb link: %s\n", getTraceLink(processor.teamSlug, processor.environmentSlug, processor.serviceName, traceId))
	}
}
func (processor localVisualizationsSpanProcessor) Shutdown(context.Context) error   { return nil }
func (processor localVisualizationsSpanProcessor) ForceFlush(context.Context) error { return nil }

func getTraceLink(teamSlug string, environmentSlug string, serviceName string, traceID string) string {
	return fmt.Sprintf("http://ui.honeycomb.io/%s/environments/%s/datasets/%s/trace?trace_id=%s", teamSlug, environmentSlug, serviceName, traceID)
}
