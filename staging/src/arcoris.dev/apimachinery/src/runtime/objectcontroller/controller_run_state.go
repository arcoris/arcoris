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

package objectcontroller

// startRun marks c running or returns ErrAlreadyRunning.
//
// The guard is intentionally per Controller. It prevents overlapping worker
// sets from consuming the same queue through one controller instance while
// still allowing a later Run after the previous one has fully returned.
func (c *Controller) startRun() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.running {
		return ErrAlreadyRunning
	}
	c.running = true
	return nil
}

// finishRun marks c available for a later sequential Run call.
func (c *Controller) finishRun() {
	c.mu.Lock()
	c.running = false
	c.mu.Unlock()
}
