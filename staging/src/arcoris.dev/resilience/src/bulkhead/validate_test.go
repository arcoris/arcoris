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
	panicassert.RequireErrorIs(t, ErrNilBulkhead, func() { _ = nilBulkhead.Snapshot() })
	panicassert.RequireErrorIs(t, ErrNilBulkhead, func() { _ = nilBulkhead.Revision() })
	panicassert.RequireErrorIs(t, ErrNilBulkhead, func() { _ = nilBulkhead.SetLimit(1) })
	panicassert.RequireErrorIs(t, ErrNilBulkhead, func() { _, _, _ = nilBulkhead.TryAcquire() })
	panicassert.RequireErrorIs(t, ErrNilBulkhead, func() { _, _, _ = nilBulkhead.TryAcquireAmount(1) })
	panicassert.RequireErrorIs(t, ErrNilBulkhead, func() { _, _, _ = nilBulkhead.TryAcquireAmount(0) })

	var zero Bulkhead
	panicassert.RequireErrorIs(t, ErrUninitializedBulkhead, func() { _ = zero.Snapshot() })
	panicassert.RequireErrorIs(t, ErrUninitializedBulkhead, func() { _ = zero.Revision() })
	panicassert.RequireErrorIs(t, ErrUninitializedBulkhead, func() { _ = zero.SetLimit(1) })
	panicassert.RequireErrorIs(t, ErrUninitializedBulkhead, func() { _, _, _ = zero.TryAcquire() })
	panicassert.RequireErrorIs(t, ErrUninitializedBulkhead, func() { _, _, _ = zero.TryAcquireAmount(1) })
	panicassert.RequireErrorIs(t, ErrUninitializedBulkhead, func() { _, _, _ = zero.TryAcquireAmount(0) })
}

func TestLeasePanicsOnNilOrUninitializedReceiver(t *testing.T) {
	t.Parallel()

	var nilLease *Lease
	panicassert.RequireErrorIs(t, ErrInvalidLease, func() { _ = nilLease.Amount() })
	panicassert.RequireErrorIs(t, ErrInvalidLease, func() { _ = nilLease.Released() })
	panicassert.RequireErrorIs(t, ErrInvalidLease, func() { _ = nilLease.Release() })
	panicassert.RequireErrorIs(t, ErrInvalidLease, func() { _, _ = nilLease.TryRelease() })

	var zero Lease
	panicassert.RequireErrorIs(t, ErrInvalidLease, func() { _ = zero.Amount() })
	panicassert.RequireErrorIs(t, ErrInvalidLease, func() { _ = zero.Released() })
	panicassert.RequireErrorIs(t, ErrInvalidLease, func() { _ = zero.Release() })
	panicassert.RequireErrorIs(t, ErrInvalidLease, func() { _, _ = zero.TryRelease() })
}
