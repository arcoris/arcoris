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

package objectwatch

// EventKind identifies the semantic shape of one watch stream item.
type EventKind uint8

const (
	// EventChanged carries one committed objectstore.Change.
	EventChanged EventKind = iota + 1
	// EventBookmark reports stream progress without a mutation.
	EventBookmark
	// EventRestartRequired terminates a stream whose continuity is lost.
	EventRestartRequired
)

// IsValid reports whether k is a supported watch event kind.
func (k EventKind) IsValid() bool {
	return k == EventChanged || k == EventBookmark || k == EventRestartRequired
}

// String returns stable lower-camel diagnostic text.
func (k EventKind) String() string {
	switch k {
	case EventChanged:
		return "changed"
	case EventBookmark:
		return "bookmark"
	case EventRestartRequired:
		return "restartRequired"
	default:
		return "unknown"
	}
}
