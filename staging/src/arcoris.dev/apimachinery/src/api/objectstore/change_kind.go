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

package objectstore

// ChangeKind identifies the committed transition shape represented by Change.
type ChangeKind uint8

const (
	// ChangeCreated reports that a key became live.
	ChangeCreated ChangeKind = iota + 1

	// ChangeUpdated reports that one live state replaced another live state.
	ChangeUpdated

	// ChangeDeleted reports that a live state was tombstoned.
	ChangeDeleted
)

// IsValid reports whether k is one of the supported committed transitions.
func (k ChangeKind) IsValid() bool {
	return k == ChangeCreated || k == ChangeUpdated || k == ChangeDeleted
}

// String returns stable lower-case diagnostic text for k.
func (k ChangeKind) String() string {
	switch k {
	case ChangeCreated:
		return "created"
	case ChangeUpdated:
		return "updated"
	case ChangeDeleted:
		return "deleted"
	default:
		return "unknown"
	}
}
