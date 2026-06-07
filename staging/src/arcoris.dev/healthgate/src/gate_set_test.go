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

package healthgate

import (
	"errors"
	"testing"

	"arcoris.dev/health"
)

func TestGateSetNormalizesEmptyResultName(t *testing.T) {
	t.Parallel()

	gate := mustGate(t)
	if err := gate.Set(health.Degraded("", health.ReasonOverloaded, "overloaded")); err != nil {
		t.Fatalf("Set() = %v, want nil", err)
	}

	result := gate.Snapshot()
	if result.Name != "ready_gate" || result.Status != health.StatusDegraded {
		t.Fatalf("Snapshot() = %+v, want named degraded result", result)
	}
}

func TestGateSetNilReceiverReturnsNilChecker(t *testing.T) {
	t.Parallel()

	var gate *Gate
	if err := gate.Healthy(); !errors.Is(err, health.ErrNilChecker) {
		t.Fatalf("Healthy(nil) = %v, want ErrNilChecker", err)
	}
}
