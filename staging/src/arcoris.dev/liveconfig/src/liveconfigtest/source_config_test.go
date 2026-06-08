// Copyright 2026 The ARCORIS Authors.
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

package liveconfigtest

import "testing"

func TestConfigSourceClonesPublishedValues(t *testing.T) {
	cfg := NewConfig()
	src := NewConfigSource(cfg)

	MutateConfig(&cfg)
	RequireConfigEqual(t, src.Current(), NewConfig())

	next := NewConfigVersion(2)
	PublishConfig(src, next)
	MutateConfig(&next)
	RequireConfigEqual(t, src.Current(), NewConfigVersion(2))
}

func TestConfigSourceStampedPublication(t *testing.T) {
	src := NewEmptyConfigSource()

	stamped := PublishConfigStamped(src, NewConfigVersion(1))

	RequireStampedNonZeroRevision(t, stamped)
	RequireConfigStampedValue(t, stamped, NewConfigVersion(1))
}
