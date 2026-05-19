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

package admission

import "arcoris.dev/value/maybe"

// Maybe is the admission-local optional-value representation.
//
// The alias keeps admission code and tests phrased in package vocabulary while
// preserving the exact representation and behavior from arcoris.dev/value/maybe.
// Result uses Maybe internally for optional grants and optional metadata.
type Maybe[T any] = maybe.Maybe[T]

// some returns an admission-local optional value containing value.
func some[T any](value T) Maybe[T] {
	return maybe.Some(value)
}

// none returns an empty admission-local optional value.
func none[T any]() Maybe[T] {
	return maybe.None[T]()
}
