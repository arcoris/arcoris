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

package capacity_test

import (
	"testing"

	"arcoris.dev/capacity"
)

func TestLedgerRevisionAdvancesOncePerCommittedMutation(t *testing.T) {
	t.Parallel()

	ledger := capacity.NewLedger(10)
	rev0 := ledger.Revision()

	reservation, snap1, ok := ledger.TryReserve(4)
	if !ok {
		t.Fatal("reservation failed")
	}
	if want := rev0.Next(); snap1.Revision != want {
		t.Fatalf("reserve revision = %d, want %d", snap1.Revision, want)
	}

	snap2 := ledger.SetLimit(12)
	if want := snap1.Revision.Next(); snap2.Revision != want {
		t.Fatalf("set limit revision = %d, want %d", snap2.Revision, want)
	}

	snap3 := reservation.Release()
	if want := snap2.Revision.Next(); snap3.Revision != want {
		t.Fatalf("release revision = %d, want %d", snap3.Revision, want)
	}
}
