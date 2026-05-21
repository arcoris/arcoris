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

package liveconfig

// EqualFunc reports whether two canonical configuration values are logically
// equivalent.
//
// When configured, Holder uses EqualFunc to suppress no-op publications. If the
// current value and candidate are equal, Apply returns Changed=false and does
// not advance the source revision. Without EqualFunc, every valid candidate is
// published.
//
// EqualFunc is evaluated after clone, normalization, and validation. It should
// compare the logical configuration state, not incidental representation details
// that the normalizer has already canonicalized.
type EqualFunc[T any] func(a, b T) bool

// equalValue reports whether cfg's equality function considers a and b equal.
//
// A nil equality function intentionally means "always changed" so callers that
// do not opt into no-op suppression still receive a fresh revision for every
// valid Apply.
func equalValue[T any](cfg config[T], a, b T) bool {
	if cfg.equal == nil {
		return false
	}

	return cfg.equal(a, b)
}
