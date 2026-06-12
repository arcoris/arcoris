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

package codec

// Target names an API document model supported by a codec.
//
// Target is closed-world in v1 because targets are framework document models,
// not codec implementation names. Adding a target means the framework has a new
// API document model, not merely that a new format implementation exists.
type Target string

const (
	// TargetValue names api/value.Value documents.
	TargetValue Target = "value"

	// TargetObject names value-backed api/object envelopes.
	TargetObject Target = "object"

	// TargetObjectOwnership names api/objectownership.State values.
	TargetObjectOwnership Target = "object_ownership"
)

// String returns the target text.
//
// Target text is stable framework metadata, not a user-facing label.
func (t Target) String() string {
	return string(t)
}

// IsZero reports whether t is absent.
//
// A zero Target is never valid codec metadata, but IsZero is intentionally only
// an absence check.
func (t Target) IsZero() bool {
	return t == ""
}
