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

package bulkhead

import (
	"errors"
	"sync"
)

// Permit is a release-once capability returned by a successful TryAcquire call.
//
// A Permit owns one unit of Limiter capacity until Release is called. Release is
// idempotent: repeated calls on the same Permit do not release more than one
// unit and cannot underflow the limiter.
type Permit struct {
	// noCopy lets go vet report accidental Permit copies after first use.
	noCopy noCopy

	// mu protects released.
	mu sync.Mutex

	// limiter owns the capacity acquired by this permit.
	limiter *Limiter

	// released reports whether this permit has already returned its capacity.
	released bool
}

var (
	// ErrReleaseUnderflow reports an impossible release against an empty limiter.
	//
	// The error is used as a panic sentinel because it indicates corruption of the
	// permit lifecycle or limiter invariants rather than an expected caller error.
	ErrReleaseUnderflow = errors.New("bulkhead: release underflow")
)

// Release returns this permit's capacity to the limiter.
//
// Release is safe to call more than once. A nil Permit is treated as an already
// released permit and is ignored; this keeps cleanup paths simple when callers
// guard on a nil permit returned from denied admission.
func (p *Permit) Release() {
	if p == nil {
		return
	}

	l, ok := p.markReleased()
	if !ok || l == nil {
		return
	}

	l.release()
}

// Released reports whether this permit has already been released.
func (p *Permit) Released() bool {
	if p == nil {
		return true
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	return p.released
}

// markReleased records the first release and returns the owning limiter.
func (p *Permit) markReleased() (*Limiter, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.released {
		return nil, false
	}

	p.released = true
	return p.limiter, true
}
