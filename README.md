[![OSS Lifecycle](https://img.shields.io/osslifecycle/honeycombio/honeycomb-opentelemetry-go)](https://github.com/honeycombio/home/blob/main/honeycomb-oss-lifecycle-and-practices.md)
[![Build Status](https://circleci.com/gh/honeycombio/honeycomb-opentelemetry-go.svg?style=shield)](https://circleci.com/gh/honeycombio/honeycomb-opentelemetry-go)

# Honeycomb OpenTelemetry Distro for Go

**STATUS: this library is BETA.**
You're welcome to try it, and let us know your feedback in the issues!

This is Honeycomb's distribution of OpenTelemetry for Go.
It makes getting started with OpenTelemetry and Honeycomb easier!

Latest release built with:

- [OpenTelemetry v1.18.0/v0.41.0](https://github.com/open-telemetry/opentelemetry-go/releases/tag/v1.18.0)

## Why would I want to use this?

- Streamlined configuration for sending data to Honeycomb!
- Easy interop with existing instrumentation with OpenTelemetry!
- Deterministic sampling!
- Multi-span attributes!
- Dynamic attributes!

## Where's most of the code?

This package is a _layer_ on top of the core package, which you can find in [here](https://github.com/honeycombio/otel-config-go). As such, this package only contains Honeycomb-specific functionality.

Our goal is to have the `otel-config-go` package be donated to the [opentelemetry-go-contrib](https://github.com/open-telemetry/opentelemetry-go-contrib) project as a blessed, vendor-neutral way to get started.

## License

[Apache 2.0 License](./LICENSE).
