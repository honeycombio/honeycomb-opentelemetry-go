# Honeycomb OpenTelemetry Distro Changelog

## 0.11.0 (2024-05-10)

### Maintenance

- maint: update example deps (#213)
- maint: NewBaggageSpanProcessor replaced by go.opentelemetry.io/contrib/processors/baggage/baggagetrace.New (#204)
  - Includes updating OTel dependencies to 1.26.0
- maint: Update README title (#207)
- maint: Update project status and add status link in README (#206)
- maint: update otel dependencies to 1.25.0 (#202)
- maint: Update ubuntu image in workflows to latest (#196)
- maint(deps): bump go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux from 0.48.0 to 0.51.0 in /examples/webhook-listener-triggers (#208)
- maint(deps): bump go.opentelemetry.io/otel/exporters/stdout/stdouttrace from 1.25.0 to 1.26.0 in the otel-dependencies group across 1 directory (#212)
- maint(deps): bump github.com/stretchr/testify from 1.8.4 to 1.9.0 (#199)
- maint(deps): bump go.opentelemetry.io/otel from 1.23.0 to 1.24.0 (#191)
- maint(deps): bump go.opentelemetry.io/otel from 1.23.0 to 1.24.0 in /examples/baggage (#192)

## v0.10.0 (2024-03-06)

### Enhancements

- feat: support Classic Ingest Keys (#193)

### Maintenance

- maint(deps): bump go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux from 0.46.1 to 0.48.0 in /examples/webhook-listener-triggers (#188)
- maint(deps): bump the otel-dependencies group with 1 update (#186)
- maint(deps): bump go.opentelemetry.io/otel from 1.21.0 to 1.23.0 in /examples/baggage (#187)
- maint(deps): bump go.opentelemetry.io/otel from 1.21.0 to 1.23.0 (#185)
- maint(deps): bump github.com/honeycombio/otel-config-go from 1.13.0 to 1.13.1 in /examples/baggage (#180)
- maint(deps): bump github.com/honeycombio/otel-config-go from 1.13.0 to 1.13.1 in /examples/webhook-listener-triggers (#179)
- maint(deps): bump github.com/honeycombio/otel-config-go from 1.13.0 to 1.13.1 (#178)
- maint: update codeowners to pipeline-team (#177)
- maint: update codeowners to pipeline (#176)

## v0.9.0 (2023-11-23)

### 💥 Breaking Changes 💥

Minimum Go Version is 1.20

### Maintenance

- maint(deps): update otel dependencies to 1.21.0 (#174)
- maint(deps): update github.com/honeycombio/otel-config-go 1.13.0 (#174)

## v0.8.1 (2023-09-22)

### Maintenance

- maint(deps): bump go.opentelemetry.io/otel/exporters/stdout/stdouttrace from 1.16.0 to 1.18.0 (#160)
- maint(deps): bump github.com/honeycombio/otel-config-go from 1.12.0 to 1.12.1 (#159)

## v0.8.0 (2023-08-16)

### 💥 Breaking Changes 💥

In previous versions, incompatible resource configurations would fail silently.
Now an error is returned so it is clear when configuration is incompatible.

### Maintenance

- maint(deps): bump github.com/honeycombio/otel-config-go from 1.11.0 to 1.12.0 (#154)
- maint(deps): bump github.com/honeycombio/otel-config-go from 1.10.0 to 1.11.0 (#152)

## v0.7.0 (2023-06-02)

### 💥 Breaking Changes 💥

Packages for the Metrics API have been moved as the API implementation has stablized in OTel Go v1.16.0.

- `go.opentelemetry.io/otel/metric/global` -> `go.opentelemetry.io/otel`
- `go.opentelemetry.io/otel/metric/instrument` -> `go.opentelemetry.io/otel/metric`

Imports of these packages in your application will need to be updated.

### Maintenance

- maint(deps): bump github.com/honeycombio/otel-config-go from 1.9.0 to 1.10.0 (#148)
  - includes: bump otel from 1.15.1/0.38.1 to 1.16.0/0.39.0

## v0.6.0 (2023-05-16)

### 💥 Breaking Changes 💥

- maint: drop go 1.18 (#144) | [@vreynolds](https://github.com/vreynolds)

### Maintenance

- maint: bump otel from 1.14.0 to 1.15.1 (#145)

## v0.5.4 (2023-04-20)

### Maintenance

- maint: Update launcher & set default endpoint (#135) | [@MikeGoldsmith](https://github.com/MikeGoldsmith)

## v0.5.3 (2023-04-12)

### Maintenance

- Update launcher & set default endpoint (#135) | [@MikeGoldsmith](https://github.com/MikeGoldsmith)
- Add go 1.20 to ci (#133) | [@vreynolds](https://github.com/vreynolds)

## v0.5.2 (2023-04-05)

### Maintenance

- maint: update launcher to 1.6.0 (#130) | [@JamieDanielson](https://github.com/jamiedanielson)

## v0.5.1 (2023-03-28)

### Maintenance

- maint: update launcher to v0.3.1 (#128) | [@JamieDanielson](https://github.com/jamiedanielson)

## v0.5.0 (2023-03-02)

### Maintenance

- chore: Rename webapp example to baggage, give it a module (#116) | [@cartermp](https://github.com/cartermp)
- maint: Add go.work to manage multiple modules (#117) | [@cartermp](https://github.com/cartermp)
- maint(deps): update otel deps to 1.14.0, launcher to 0.3.0 (#125) | [@JamieDanielson](https://github.com/jamiedanielson)
- maint(deps): bump github.com/stretchr/testify from 1.8.1 to 1.8.2 (#123)
- build(deps): bump golang.org/x/net from 0.4.0 to 0.7.0 in /examples/webhook-listener-triggers (#120)
- build(deps): bump golang.org/x/net from 0.4.0 to 0.7.0 in /examples/baggage (#119)
- maint(deps): bump golang.org/x/net from 0.4.0 to 0.7.0 (#118)

## v0.4.2 (2023-02-01)

### Maintenance

- chore: update to latest launcher and otel pkgs (#114) | [@cartermp](https://github.com/cartermp)

## v0.4.1 (2023-01-19)

### Fixes

- Use configured logger to print messages (#103) | [@martin308](https://github.com/martin308)

### Maintenance

- Add smoke test to circle (#107) | [@JamieDanielson](https://github.com/jamiedanielson)
- Add smoke tests for traces (#102) | [@JamieDanielson](https://github.com/jamiedanielson)
- Update launcher to use new repo (#109) | [@MikeGoldsmith](https://github.com/MikeGoldsmith)

## v0.4.0 (2023-01-04)

### Fixes

- Don't error on misconfiguration; warn instead (#98) | [@cartermp](https://github.com/cartermp)

### Maintenance

- Update Launcher and OTel packages (#100) | [@MikeGoldsmith](https://github.com/MikeGoldsmith)
- Fix typo in test file (#99) | [@JamieDanielson](https://github.com/jamiedanielson)
- Update readme with latest otel version (#88) | [@vreynolds](https://github.com/vreynolds)
- Update validate PR title workflow (#91) | [@pkanal](https://github.com/pkanal)
- Validate PR title (#90) | [@pkanal](https://github.com/pkanal)
- Give dependabot PRs a better title (#93) | [@kentquirk](https://github.com/kentquirk)

## v0.3.0 (2022-10-31)

### Changes

- Minimum required Go version is 1.18 (#84)

### Maintenance

- Remove timestamp field from example trigger hook (#81) | @passcod
- Update launcher to latest (#80, #86) | @MikeGoldsmith @vreynolds
  - fix unconditional debug statements
  - update OTEL packages
- Bump go.opentelemetry.io/otel/exporters/stdout/stdouttrace from 1.9.0 to 1.11.1 (#84)
- Bump go.opentelemetry.io/otel/sdk from 1.9.0 to 1.10.0 (#76)

## v0.2.0 (2022-08-24)

### Enhancements

- Add local visualizations exporter (#66) | [@cartermp](https://github.com/cartermp)
- Add support for separate traces and metrics API key and dataset (#72) | [@MikeGoldsmith](https://github.com/MikeGoldsmith)
- Disable metrics by default (#70) | [@MikeGoldsmith](https://github.com/MikeGoldsmith)
- Add support for Honeycomb endpoint environment variables (#65) | [@MikeGoldsmith](https://github.com/MikeGoldsmith)
- Add OTLP version header (#64) | [@vreynolds](https://github.com/vreynolds)

### Maintenance

- Add webhook triggers example (#68) | [@vreynolds](https://github.com/vreynolds)
- Add test matrix and nightly (#67) | [@vreynolds](https://github.com/vreynolds)

## v0.1.2 (2022-08-17)

## Fixed

- Set base exporter endpoint when setting up vendor opts (#56) | [@MikeGoldsmith](https://github.com/MikeGoldsmith)
- Set log level to debug when debug option is set (#55) | [@MikeGoldsmith](https://github.com/MikeGoldsmith)

### Maintenance

- Add baggage processor tests (#58) | [@MikeGoldsmith](https://github.com/MikeGoldsmith)
- Add missing license headers (#57) | [@MikeGoldsmith](https://github.com/MikeGoldsmith)
- More descriptive errors (#60) | [@cartermp](https://github.com/cartermp)

## v0.1.1 (2022-08-12)

### Fixed

- Update module path to match repo path (#46) | [@MikeGoldsmith](https://github.com/MikeGoldsmith)

### Maintenance

- Update README to clarify where most of the code lives (#45) | [@cartermp](https://github.com/cartermp)

## v0.1.0 (2022-08-12)

Initial release
