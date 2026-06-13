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

// ListItem is one committed live object returned by Store.List.
type ListItem struct {
	// Key is the authoritative storage identity for State.
	Key Key

	// State is the detached committed live state for Key.
	State State
}

// Clone returns a detached copy of i.
func (i ListItem) Clone() ListItem {
	i.State = i.State.Clone()

	return i
}
