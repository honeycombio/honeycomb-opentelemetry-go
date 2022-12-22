package honeycomb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSpanLinkExporterBuildsValidUrl(t *testing.T) {
	t.Run("classic", func(t *testing.T) {
		assert.Equal(t, "https://ui.honeycomb.io/my-team/datasets/my-service/trace?trace_id", buildTraceLinkUrl(true, "my-team", "my-env", "my-service"))
	})
	t.Run("environment", func(t *testing.T) {
		assert.Equal(t, "https://ui.honeycomb.io/my-team/environments/my-env/datasets/my-service/trace?trace_id", buildTraceLinkUrl(false, "my-team", "my-env", "my-service"))
	})
}
