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
	"arcoris.dev/snapshot"
)

func requireSnapshotValue(t *testing.T, snap snapshot.Snapshot[capacity.Snapshot], limit, reserved, available, debt capacity.Amount) {
	t.Helper()

	if snap.Revision.IsZero() {
		t.Fatalf("snapshot revision is zero")
	}
	if !snap.Value.IsValid() {
		t.Fatalf("snapshot value is invalid: %+v", snap.Value)
	}
	if snap.Value.Limit != limit || snap.Value.Reserved != reserved || snap.Value.Available != available || snap.Value.Debt != debt {
		t.Fatalf(
			"snapshot value = {Limit:%d Reserved:%d Available:%d Debt:%d}, want {Limit:%d Reserved:%d Available:%d Debt:%d}",
			snap.Value.Limit,
			snap.Value.Reserved,
			snap.Value.Available,
			snap.Value.Debt,
			limit,
			reserved,
			available,
			debt,
		)
	}
}
