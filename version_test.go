package honeycomb_test

import (
	"regexp"
	"testing"

	"github.com/honeycombio/otel/honeycomb"
	"github.com/stretchr/testify/assert"
)

var versionRegex = regexp.MustCompile(`^v?([0-9]+)(\.[0-9]+)?(\.[0-9]+)?` +
	`(-([0-9A-Za-z\-]+(\.[0-9A-Za-z\-]+)*))?` +
	`(\+([0-9A-Za-z\-]+(\.[0-9A-Za-z\-]+)*))?$`)

func TestVersionSemver(t *testing.T) {
	v := honeycomb.Version
	assert.NotNil(t, versionRegex.FindStringSubmatch(v), "version is not semver: %s", v)
}
