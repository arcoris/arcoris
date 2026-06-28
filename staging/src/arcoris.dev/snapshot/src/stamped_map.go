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

package snapshot

// MapStamped transforms a stamped snapshot value while preserving Revision and
// Updated.
//
// MapStamped does not clone either the input or output value. The source that
// produced the stamped snapshot owns any immutability or clone guarantees.
// MapStamped panics when fn is nil because a nil mapping function is a
// programmer error.
func MapStamped[R comparable, A any, B any](s Stamped[R, A], fn func(A) B) Stamped[R, B] {
	if fn == nil {
		panic("snapshot: nil stamped map function")
	}

	return Stamped[R, B]{
		Revision: s.Revision,
		Updated:  s.Updated,
		Value:    fn(s.Value),
	}
}
