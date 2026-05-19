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

import "testing"

func TestBulkheadPanicsOnNilOrUninitializedReceiver(t *testing.T) {
	t.Parallel()

	var nilBulkhead *Bulkhead
	requirePanic(t, errNilBulkhead, func() { _ = nilBulkhead.Snapshot() })
	requirePanic(t, errNilBulkhead, func() { _ = nilBulkhead.Revision() })
	requirePanic(t, errNilBulkhead, func() { _ = nilBulkhead.SetLimit(1) })
	requirePanic(t, errNilBulkhead, func() { _, _, _ = nilBulkhead.TryAcquire() })

	var zero Bulkhead
	requirePanic(t, errUninitializedBulkhead, func() { _ = zero.Snapshot() })
	requirePanic(t, errUninitializedBulkhead, func() { _ = zero.Revision() })
	requirePanic(t, errUninitializedBulkhead, func() { _ = zero.SetLimit(1) })
	requirePanic(t, errUninitializedBulkhead, func() { _, _, _ = zero.TryAcquire() })
}

func TestLeasePanicsOnNilOrUninitializedReceiver(t *testing.T) {
	t.Parallel()

	var nilLease *Lease
	requirePanic(t, errInvalidLease, func() { _ = nilLease.Amount() })
	requirePanic(t, errInvalidLease, func() { _ = nilLease.Released() })
	requirePanic(t, errInvalidLease, func() { _ = nilLease.Release() })
	requirePanic(t, errInvalidLease, func() { _, _ = nilLease.TryRelease() })

	var zero Lease
	requirePanic(t, errInvalidLease, func() { _ = zero.Amount() })
	requirePanic(t, errInvalidLease, func() { _ = zero.Released() })
	requirePanic(t, errInvalidLease, func() { _ = zero.Release() })
	requirePanic(t, errInvalidLease, func() { _, _ = zero.TryRelease() })
}
