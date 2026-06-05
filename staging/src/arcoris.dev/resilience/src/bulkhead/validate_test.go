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

import (
	"testing"

	panicassert "arcoris.dev/testutil/panic"
)

func TestBulkheadPanicsOnNilOrUninitializedReceiver(t *testing.T) {
	t.Parallel()

	var nilBulkhead *Bulkhead
	panicassert.RequireMessage(t, errNilBulkhead, func() { _ = nilBulkhead.Snapshot() })
	panicassert.RequireMessage(t, errNilBulkhead, func() { _ = nilBulkhead.Revision() })
	panicassert.RequireMessage(t, errNilBulkhead, func() { _ = nilBulkhead.SetLimit(1) })
	panicassert.RequireMessage(t, errNilBulkhead, func() { _, _, _ = nilBulkhead.TryAcquire() })
	panicassert.RequireMessage(t, errNilBulkhead, func() { _, _, _ = nilBulkhead.TryAcquireAmount(1) })
	panicassert.RequireMessage(t, errNilBulkhead, func() { _, _, _ = nilBulkhead.TryAcquireAmount(0) })
	panicassert.RequireMessage(t, errNilBulkhead, func() { _ = nilBulkhead.TryAdmit(Request{Amount: 1}) })
	panicassert.RequireMessage(t, errNilBulkhead, func() { _ = nilBulkhead.TryAdmit(Request{Amount: 0}) })

	var zero Bulkhead
	panicassert.RequireMessage(t, errUninitializedBulkhead, func() { _ = zero.Snapshot() })
	panicassert.RequireMessage(t, errUninitializedBulkhead, func() { _ = zero.Revision() })
	panicassert.RequireMessage(t, errUninitializedBulkhead, func() { _ = zero.SetLimit(1) })
	panicassert.RequireMessage(t, errUninitializedBulkhead, func() { _, _, _ = zero.TryAcquire() })
	panicassert.RequireMessage(t, errUninitializedBulkhead, func() { _, _, _ = zero.TryAcquireAmount(1) })
	panicassert.RequireMessage(t, errUninitializedBulkhead, func() { _, _, _ = zero.TryAcquireAmount(0) })
	panicassert.RequireMessage(t, errUninitializedBulkhead, func() { _ = zero.TryAdmit(Request{Amount: 1}) })
	panicassert.RequireMessage(t, errUninitializedBulkhead, func() { _ = zero.TryAdmit(Request{Amount: 0}) })
}

func TestLeasePanicsOnNilOrUninitializedReceiver(t *testing.T) {
	t.Parallel()

	var nilLease *Lease
	panicassert.RequireMessage(t, errInvalidLease, func() { _ = nilLease.Amount() })
	panicassert.RequireMessage(t, errInvalidLease, func() { _ = nilLease.Released() })
	panicassert.RequireMessage(t, errInvalidLease, func() { _ = nilLease.Release() })
	panicassert.RequireMessage(t, errInvalidLease, func() { _, _ = nilLease.TryRelease() })

	var zero Lease
	panicassert.RequireMessage(t, errInvalidLease, func() { _ = zero.Amount() })
	panicassert.RequireMessage(t, errInvalidLease, func() { _ = zero.Released() })
	panicassert.RequireMessage(t, errInvalidLease, func() { _ = zero.Release() })
	panicassert.RequireMessage(t, errInvalidLease, func() { _, _ = zero.TryRelease() })
}
