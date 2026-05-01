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

package lifecycle

// Done returns a channel that is closed when the lifecycle reaches a terminal
// state.
//
// Done is a signal-only channel. It does not report whether the terminal state is
// StateStopped or StateFailed. Call Snapshot or WaitTerminal to inspect the final
// state and failure cause.
//
// The returned channel is stable for the lifetime of the Controller. It is closed
// at most once.
func (c *Controller) Done() <-chan struct{} {
	return c.doneSignal()
}
