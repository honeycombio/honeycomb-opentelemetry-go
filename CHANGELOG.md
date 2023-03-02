# Honeycomb OpenTelemetry Distro Changelog

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
