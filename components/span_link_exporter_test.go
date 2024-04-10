// Copyright Honeycomb Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package components

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
