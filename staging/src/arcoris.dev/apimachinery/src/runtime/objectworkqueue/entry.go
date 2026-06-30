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

package objectworkqueue

// entry is the mutable state for one tracked object key.
type entry struct {
	// item is the original public work item associated with this entry.
	item Item

	// state records whether item is queued or currently processing.
	state entryState

	// dirty records that Add was called while item was processing.
	dirty bool

	// prev points to the previous queued entry while this entry is queued.
	prev *entry

	// next points to the next queued entry while this entry is queued.
	next *entry
}

// entryState is the private tracked-item state machine.
type entryState uint8

const (
	// stateQueued means the item is waiting in FIFO order for Get.
	stateQueued entryState = iota + 1

	// stateProcessing means the item has been returned by Get and awaits Done.
	stateProcessing
)
