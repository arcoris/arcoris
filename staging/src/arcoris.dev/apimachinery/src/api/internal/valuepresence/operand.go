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

package valuepresence

import "arcoris.dev/apimachinery/api/value"

// Operand carries one traversal-side value while preserving presence.
//
// Presence is intentionally separate from value.Value because explicit null,
// scalar zeroes, and empty composites are real payload data. Only an absent
// Operand means the field, key, item, or selector was not found.
type Operand struct {
	value   value.Value
	present bool
}

// Present wraps a concrete payload value as present traversal input.
func Present(v value.Value) Operand {
	return Operand{
		value:   v,
		present: true,
	}
}

// Absent returns the canonical absent traversal input.
func Absent() Operand {
	return Operand{}
}

// From converts a conventional value, ok lookup result into an Operand.
func From(v value.Value, ok bool) Operand {
	if !ok {
		return Absent()
	}

	return Present(v)
}

// Clone returns a detached copy of o.
//
// Absent operands stay canonical absent operands. Present operands clone their
// payload so traversal packages can preserve immutable-by-convention behavior
// without knowing value.Value internals.
func (o Operand) Clone() Operand {
	if !o.present {
		return Absent()
	}

	return Present(o.value.Clone())
}

// Present reports whether o carries a concrete payload value.
func (o Operand) Present() bool {
	return o.present
}

// Absent reports whether o represents a missing field, key, item, or selector.
func (o Operand) Absent() bool {
	return !o.present
}

// Value returns o's payload value.
//
// For absent operands this is the zero value.Value. Callers that need to branch
// on presence should use ValueOK.
func (o Operand) Value() value.Value {
	return o.value
}

// ValueOK returns o's payload value and presence flag together.
func (o Operand) ValueOK() (value.Value, bool) {
	return o.value, o.present
}
