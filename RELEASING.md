# Creating a new release

1. Update the version string in `version.go`.
1. If this release includes an upgrade to the underlying OTLP proto version, update `otlpProtoVersionValue` in `honeycomb.go` (run `go list -m all` to see which version was selected for `go.opentelemetry.io/proto/otlp`)
1. If this release updates OTel versions, update "OTel version this is built with" in the `README.md`
1 Update `CHANGELOG.md` with the changes since the last release. Consider automating with a command such as these two:
  a. `git log $(git describe --tags --abbrev=0)..HEAD --no-merges --oneline > new-in-this-release.log`
  a. `git log --pretty='%C(green)%d%Creset- %s | [%an](https://github.com/)'`
1. Commit changes, push, and open a release preparation pull request for review.
1. Once the pull request is merged, fetch the updated `main` branch.
1. Apply a tag for the new version on the merged commit (e.g. `git tag -a v1.3.0 -m "v1.3.0"`)
1. Push the tag upstream (this will kick off the release pipeline in CI) e.g. `git push origin v1.3.0`
1. Ensure that there is a draft GitHub release created as part of CI publish steps
1. Click "generate release notes" in Github for full changelog notes and any new contributors
1. Click the prerelease checkbox and publish the Github draft release
