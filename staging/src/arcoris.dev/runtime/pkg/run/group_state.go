/*
  Copyright 2026 The ARCORIS Authors

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
*/

package run

// reserveTask reserves name, assigns a submission sequence, and increments the
// internal WaitGroup before the goroutine is started.
//
// Wait closes the group under the same mutex. Keeping the closed check, name
// registration, sequence assignment, and WaitGroup.Add under one lock prevents
// Go from racing with Wait.
func (g *Group) reserveTask(name string) uint64 {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.closed {
		panic(errGroupClosed)
	}
	if _, exists := g.names[name]; exists {
		panic(errDuplicateTaskName)
	}

	seq := g.nextSeq
	g.nextSeq++
	g.names[name] = struct{}{}
	g.wg.Add(1)

	return seq
}

// close prevents future task submissions only.
//
// Closing is deliberately smaller than cancellation or joining: it does not
// cancel the context, wait for running tasks, or build the cached Wait error.
// Those responsibilities stay with Cancel, task-error handling, and Wait.
func (g *Group) close() {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.closed = true
}
