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

package memory

import "arcoris.dev/apimachinery/api/objectstore"

// liveRecord builds an immutable live record with detached state.
func liveRecord(state objectstore.State, revision objectstore.Revision) *record {
	return &record{state: objectstore.AssignRevision(state, revision)}
}

// tombstoneRecord builds an immutable tombstone record.
func tombstoneRecord(previous objectstore.State, revision objectstore.Revision) *record {
	return &record{
		state:          previous.Clone(),
		deleteRevision: revision,
		deleted:        true,
	}
}
