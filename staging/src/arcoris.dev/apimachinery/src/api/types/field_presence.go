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

package types

// FieldPresence describes whether an object field must be present.
//
// Presence is a field-level property. It is intentionally separate from
// Type.Nullable: a required field may still allow a null value, and an optional
// field may still reject null when it is present.
type FieldPresence uint8

const (
	// PresenceUnspecified is the zero value and is invalid for finalized fields.
	PresenceUnspecified FieldPresence = iota
	// PresenceRequired means an object value must include the field key.
	PresenceRequired
	// PresenceOptional means an object value may omit the field key.
	PresenceOptional
)

// IsValid reports whether p is a usable finalized field presence.
func (p FieldPresence) IsValid() bool {
	return p == PresenceRequired || p == PresenceOptional
}

// String returns the stable diagnostic spelling of p.
func (p FieldPresence) String() string {
	switch p {
	case PresenceUnspecified:
		return "unspecified"
	case PresenceRequired:
		return "required"
	case PresenceOptional:
		return "optional"
	default:
		return "unknown"
	}
}
