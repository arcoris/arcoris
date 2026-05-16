/*
  Copyright 2026 The ARCORIS Authors

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
*/

package liveconfigtest

import "arcoris.dev/snapshot"

// NewConfigSource creates a ControlledSource for the Config fixture.
//
// Config values contain maps and slices, so this helper clones initial before
// publication. Tests that intentionally need Publisher-style no-clone behavior
// should use NewControlledSource directly.
func NewConfigSource(initial Config, opts ...snapshot.Option) *ControlledSource[Config] {
	return NewControlledSource(CloneConfig(initial), opts...)
}

// NewEmptyConfigSource creates an empty ControlledSource for Config values.
func NewEmptyConfigSource(opts ...snapshot.Option) *ControlledSource[Config] {
	return NewEmptyControlledSource[Config](opts...)
}

// PublishConfig clones cfg and publishes it to src.
func PublishConfig(src *ControlledSource[Config], cfg Config) snapshot.Snapshot[Config] {
	return src.Publish(CloneConfig(cfg))
}

// PublishConfigStamped clones cfg and publishes it to src with timestamp
// metadata.
func PublishConfigStamped(src *ControlledSource[Config], cfg Config) snapshot.Stamped[Config] {
	return src.PublishStamped(CloneConfig(cfg))
}
