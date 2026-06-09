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

package objectlifecycle

// Effect describes the visible state transition caused by a successful operation.
type Effect uint8

const (
	// EffectNone is the zero no-effect value used by zero Results.
	EffectNone Effect = iota
	// EffectFound reports that Get found committed live state.
	EffectFound
	// EffectCreated reports that a new live state was committed.
	EffectCreated
	// EffectUpdated reports that existing live state was replaced.
	EffectUpdated
	// EffectUnchanged is reserved for future no-op apply detection.
	EffectUnchanged
	// EffectDeleted reports that live state was deleted.
	EffectDeleted
)

// IsValid reports whether e is a known lifecycle effect.
func (e Effect) IsValid() bool {
	return e <= EffectDeleted
}

// String returns stable diagnostic text for e.
func (e Effect) String() string {
	switch e {
	case EffectNone:
		return "none"
	case EffectFound:
		return "found"
	case EffectCreated:
		return "created"
	case EffectUpdated:
		return "updated"
	case EffectUnchanged:
		return "unchanged"
	case EffectDeleted:
		return "deleted"
	default:
		return "unknown"
	}
}
