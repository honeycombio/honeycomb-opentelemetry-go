# Creating a new release

1. Update the version string in `version.go`.
1. If this release includes an upgrade to the underlying OTLP proto version, update `otlpProtoVersionValue` in `honeycomb.go` (run `go list -m all` to see which version was selected for `go.opentelemetry.io/proto/otlp`)
1. If this release updates OTel versions, update "OTel version this is built with" in the README.md
1. Add new release notes to the Changelog.
1. Open a PR with above changes.
1. Once the above PR is merged, tag `main` with the new version, e.g. `v0.1.0`, and push the tags. This will kick off a CI workflow, which will publish a draft GitHub release.
1. Update Release Notes on the new draft GitHub release using the Changelog entry, and publish that as Pre-release.
