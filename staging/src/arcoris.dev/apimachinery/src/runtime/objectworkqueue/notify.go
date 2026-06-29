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

// signalNotEmptyLocked wakes waiters observing q.notEmpty.
func (q *Queue) signalNotEmptyLocked() {
	close(q.notEmpty)
	q.notEmpty = make(chan struct{})
}

// signalNotFullLocked wakes waiters observing q.notFull.
func (q *Queue) signalNotFullLocked() {
	close(q.notFull)
	q.notFull = make(chan struct{})
}

// signalAllLocked wakes both item and capacity waiters.
func (q *Queue) signalAllLocked() {
	q.signalNotEmptyLocked()
	q.signalNotFullLocked()
}
