#!/usr/bin/env bats

load test_helpers/utilities

CONTAINER_NAME="app-sdk-http"
TRACER_NAME="hello-world-tracer"
METER_NAME="hello-world-meter"

setup_file() {
	echo "# 🚧" >&3
	docker-compose up --build --detach collector ${CONTAINER_NAME}
	wait_for_ready_app ${CONTAINER_NAME}
	curl --silent "http://localhost:8090"
	wait_for_traces
  # wait_for_metrics 15
}

teardown_file() {
	cp collector/data.json collector/data-results/data-${CONTAINER_NAME}.json
	docker-compose stop ${CONTAINER_NAME}
	docker-compose restart collector
	wait_for_flush
}

# TESTS

@test "Auto instrumentation produces an http request span" {
  result=$(span_names_for "go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux")
  assert_equal "$result" '"/"'
}

# @test "Manual instrumentation produces span with name of span" {
# 	result=$(span_names_for ${TRACER_NAME})
# 	assert_equal "$result" '"sleep"'
# }

# @test "Manual instrumentation adds custom attribute" {
# 	result=$(span_attributes_for ${TRACER_NAME} | jq "select(.key == \"delay_ms\").value.intValue")
# 	assert_equal "$result" '"100"'
# }

# @test "BaggageSpanProcessor: key-values added to baggage appear on child spans" {
# 	result=$(span_attributes_for ${TRACER_NAME} | jq "select(.key == \"for_the_children\").value.stringValue")
# 	assert_equal "$result" '"another important value"'
# }

# @test "Manual instrumentation produces metrics" {
#     result=$(metric_names_for ${METER_NAME})
#     assert_equal "$result" '"sheep"'
# }