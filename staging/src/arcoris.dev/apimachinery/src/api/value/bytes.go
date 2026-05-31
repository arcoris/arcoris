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

package value

// BytesValue constructs a KindBytes Value by cloning input.
//
// The caller may safely mutate v after construction. Nil input is normalized to
// an owned empty slice so byte payloads have stable non-nil storage internally.
func BytesValue(v []byte) Value {
	return Value{kind: KindBytes, bytesValue: cloneBytes(v)}
}

// cloneBytes returns an owned copy and normalizes nil to an empty slice.
//
// This helper is used by both constructors and accessors so byte payloads never
// share mutable backing arrays across the public API boundary.
func cloneBytes(v []byte) []byte {
	cloned := make([]byte, len(v))
	copy(cloned, v)

	return cloned
}
