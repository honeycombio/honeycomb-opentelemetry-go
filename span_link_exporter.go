package honeycomb

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

var zeroTime time.Time

var _ trace.SpanExporter = &Exporter{}

// New creates an Exporter with the passed options.
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

func getTraceLink(teamSlug string, environmentSlug string, serviceName string, traceID string) string {
	return fmt.Sprintf("http://ui.honeycomb.io/%s/environments/%s/datasets/%s/trace?trace_id=%s", teamSlug, environmentSlug, serviceName, traceID)
}

func (e *Exporter) ExportSpans(ctx context.Context, spans []trace.ReadOnlySpan) error {
	if len(spans) == 0 {
		return nil
	}

	stubs := tracetest.SpanStubsFromReadOnlySpans(spans)

	for i := range stubs {
		stub := &stubs[i]

		if !stub.Parent.SpanID().IsValid() {
			fmt.Printf("Trace for %s\n", stub.Name)
			traceId := stub.SpanContext.TraceID().String()
			link := getTraceLink(e.teamSlug, e.environmentSlug, e.serviceName, traceId)
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
