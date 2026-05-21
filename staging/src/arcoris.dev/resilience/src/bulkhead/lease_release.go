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


package bulkhead

import "arcoris.dev/snapshot"

// Release returns l's in-flight slot to the owning Bulkhead.
//
// Release panics if l has already been released. Double release is an ownership
// bug because it means the caller lost track of the live protected section. Use
// TryRelease when idempotent cleanup is required.
func (l *Lease) Release() snapshot.Snapshot[Snapshot] {
	l.requireReady()
	return l.reservation.Release()
}

// TryRelease returns l's in-flight slot if it is still live.
//
// On first release, TryRelease returns the resulting snapshot with true. On
// later calls, it leaves capacity unchanged and returns the current snapshot with
// false.
func (l *Lease) TryRelease() (snapshot.Snapshot[Snapshot], bool) {
	l.requireReady()
	return l.reservation.TryRelease()
}
