[![OSS Lifecycle](https://img.shields.io/osslifecycle/honeycombio/honeycomb-opentelemetry-go)](https://github.com/honeycombio/home/blob/main/honeycomb-oss-lifecycle-and-practices.md)
[![Build Status](https://circleci.com/gh/honeycombio/honeycomb-opentelemetry-go.svg?style=shield)](https://circleci.com/gh/honeycombio/honeycomb-opentelemetry-go)

# Honeycomb OpenTelemetry Distro for Go

**STATUS: this library is BETA.**
You're welcome to try it, and let us know your feedback in the issues!

This is Honeycomb's distribution of OpenTelemetry for Go.
It makes getting started with OpenTelemetry and Honeycomb easier!

Latest release built with:

- [OpenTelemetry](https://github.com/open-telemetry/opentelemetry-go/releases/tag/v1.9.0) version v1.9.0

## Why would I want to use this?

- Streamlined configuration for sending data to Honeycomb!
- Easy interop with existing instrumentation with OpenTelemetry!
- Deterministic sampling!
- Multi-span attributes!
- Dynamic attributes!

## Where's most of the code?

This package is a _layer_ on top of the core package, which you can find in our [fork of the opentelemetry-go-contrib repo](https://github.com/honeycombio/opentelemetry-go-contrib/tree/launcher/launcher). As such, it only containts Honeycomb-specific functionality.

Our immedate goal is that our fork lives upstream in the opentelemetry-go-contrib project as a blessed, vendor-neutral way to get started.

## License

[Apache 2.0 License](./LICENSE).
