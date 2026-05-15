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

package retrybudget

import (
	"math"
	"testing"
	"time"

	"arcoris.dev/snapshot"
)

func validSnapshotValue() Snapshot {
	return Snapshot{
		Kind: KindFixedWindow,
		Attempts: AttemptsSnapshot{
			Original: 10,
			Retry:    2,
		},
		Capacity: CapacitySnapshot{
			Allowed:   4,
			Available: 2,
			Exhausted: false,
		},
		Window: WindowSnapshot{
			StartedAt: time.Unix(100, 0).UTC(),
			EndsAt:    time.Unix(160, 0).UTC(),
			Duration:  time.Minute,
			Bounded:   true,
		},
		Policy: PolicySnapshot{
			Ratio:   0.2,
			Minimum: 2,
			Bounded: true,
		},
	}
}

func validGenericSnapshot() snapshot.Snapshot[Snapshot] {
	return snapshot.Snapshot[Snapshot]{
		Revision: snapshot.ZeroRevision.Next(),
		Value:    validSnapshotValue(),
	}
}

func exhaustedSnapshotValue() Snapshot {
	val := validSnapshotValue()
	val.Capacity = CapacitySnapshot{Allowed: 4, Available: 0, Exhausted: true}
	return val
}

func maxedCapacity() CapacitySnapshot {
	return CapacitySnapshot{Allowed: math.MaxUint64, Available: math.MaxUint64, Exhausted: false}
}

func requireValid(t *testing.T, valid bool) {
	t.Helper()
	if !valid {
		t.Fatal("value is invalid")
	}
}
