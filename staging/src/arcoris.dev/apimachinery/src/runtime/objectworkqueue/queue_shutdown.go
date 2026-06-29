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

// ShutDown prevents new work and wakes blocked Add and Get callers.
//
// ShutDown is idempotent. Get may continue draining already queued items, and
// Done remains valid for already processing items.
func (q *Queue) ShutDown() {
	if q == nil {
		return
	}

	q.mu.Lock()
	defer q.mu.Unlock()

	if q.shutDown {
		return
	}
	q.shutDown = true
	q.signalAllLocked()
}
