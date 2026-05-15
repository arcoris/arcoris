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

package atomicx

import "sync/atomic"

// PaddedPointer is a padded raw atomic pointer cell.
//
// PaddedPointer wraps sync/atomic.Pointer with the same layout policy used by
// the other raw padded atomicx primitives. It is intended for a small number of
// hot pointer-publication cells such as current immutable routing tables,
// runtime policies, observer lists, strategy objects, or copy-on-write read
// models.
//
// PaddedPointer does not own the pointed value. It does not clone, freeze,
// validate, retain, release, or otherwise manage the lifetime of the pointed
// object. Callers must ensure that stored values remain valid while readers can
// observe them and that published objects are not mutated unless the surrounding
// ownership protocol explicitly permits that mutation.
//
// PaddedPointer is not a snapshot source, not a cache, not an RCU
// implementation, not a hazard-pointer system, and not a replacement for
// snapshot.Publisher. Use snapshot.Publisher when callers need revisioned
// publication of immutable read models.
//
// PaddedPointer is zero-value usable.
//
// PaddedPointer must not be copied after first use. Copying a live pointer cell
// splits one logical publication point into independent cells and can corrupt
// ownership and visibility assumptions.
type PaddedPointer[T any] struct {
	noCopy noCopy
	_      CacheLinePad
	value  atomic.Pointer[T]
	_      [CacheLinePadSize - atomicPointerSize]byte
}

// Load atomically returns the pointer stored in p.
//
// Load observes only the pointer value. It does not make the pointed object
// immutable and does not provide a component-level snapshot of related fields.
func (p *PaddedPointer[T]) Load() *T {
	return p.value.Load()
}

// Store atomically stores ptr in p.
//
// Store publishes the pointer value only. The caller owns all lifetime and
// immutability rules for the pointed object.
func (p *PaddedPointer[T]) Store(ptr *T) {
	p.value.Store(ptr)
}

// Swap atomically stores new in p and returns the previous pointer.
//
// Swap is useful for owner-controlled pointer handoff. It does not validate the
// pointed value and does not manage the lifetime of the previous pointer.
func (p *PaddedPointer[T]) Swap(new *T) *T {
	return p.value.Swap(new)
}

// CompareAndSwap atomically replaces old with new when p currently stores old.
//
// CompareAndSwap is intended for advanced internal protocols where the caller
// owns the expected-pointer transition rules.
func (p *PaddedPointer[T]) CompareAndSwap(old, new *T) bool {
	return p.value.CompareAndSwap(old, new)
}
