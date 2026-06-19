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

package objectreflector

// beginRun marks r active if no other Run call is active.
//
// The guard rejects overlapping synchronization loops while still allowing a
// caller to run the same Reflector again after a previous Run returned.
func (r *Reflector) beginRun() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.running {
		return ErrAlreadyRunning
	}
	r.running = true

	return nil
}

// endRun releases the single-run guard.
func (r *Reflector) endRun() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.running = false
}
