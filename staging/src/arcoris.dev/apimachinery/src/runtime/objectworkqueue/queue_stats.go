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

// Stats is an instantaneous diagnostic observation of queue state.
//
// Stats must not be used as a synchronization mechanism.
type Stats struct {
	// Queued is the number of items waiting to be returned by Get.
	Queued int

	// Processing is the number of items returned by Get and not completed by
	// Done.
	Processing int

	// Capacity is the maximum number of queued plus processing items.
	Capacity int

	// ShutDown reports whether ShutDown has been called.
	ShutDown bool
}

// Stats returns an instantaneous diagnostic observation of queue state.
//
// A nil Queue returns the zero Stats value.
func (q *Queue) Stats() Stats {
	if q == nil {
		return Stats{}
	}

	q.mu.Lock()
	defer q.mu.Unlock()

	return Stats{
		Queued:     q.queuedLocked(),
		Processing: q.processingLocked(),
		Capacity:   q.capacity,
		ShutDown:   q.shutDown,
	}
}
