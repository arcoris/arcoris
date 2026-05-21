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
	"testing"

	"arcoris.dev/snapshot"
)

func TestDecisionStateHelpers(t *testing.T) {
	decision := Decision{Allowed: true, Reason: ReasonAllowed, Snapshot: validGenericSnapshot()}
	if !decision.IsAllowed() {
		t.Fatal("IsAllowed returned false")
	}
	if decision.IsDenied() {
		t.Fatal("IsDenied returned true")
	}

	decision.Allowed = false
	decision.Reason = ReasonExhausted
	decision.Snapshot = snapshot.Snapshot[Snapshot]{Revision: snapshot.ZeroRevision.Next(), Value: exhaustedSnapshotValue()}
	if decision.IsAllowed() {
		t.Fatal("IsAllowed returned true")
	}
	if !decision.IsDenied() {
		t.Fatal("IsDenied returned false")
	}
}

func TestDecisionIsValid(t *testing.T) {
	exhausted := snapshot.Snapshot[Snapshot]{Revision: snapshot.ZeroRevision.Next(), Value: exhaustedSnapshotValue()}
	invalidSnapshot := snapshot.Snapshot[Snapshot]{Revision: snapshot.ZeroRevision.Next(), Value: Snapshot{}}
	zeroRevision := snapshot.Snapshot[Snapshot]{Value: validSnapshotValue()}
	tests := []struct {
		name string
		val  Decision
		want bool
	}{
		{name: "zero", val: Decision{}, want: false},
		{name: "allowed", val: Decision{Allowed: true, Reason: ReasonAllowed, Snapshot: validGenericSnapshot()}, want: true},
		{name: "denied exhausted", val: Decision{Allowed: false, Reason: ReasonExhausted, Snapshot: exhausted}, want: true},
		{name: "allowed with exhausted reason", val: Decision{Allowed: true, Reason: ReasonExhausted, Snapshot: validGenericSnapshot()}, want: false},
		{name: "denied with allowed reason", val: Decision{Allowed: false, Reason: ReasonAllowed, Snapshot: validGenericSnapshot()}, want: false},
		{name: "invalid reason", val: Decision{Allowed: false, Reason: ReasonUnknown, Snapshot: validGenericSnapshot()}, want: false},
		{name: "invalid snapshot value", val: Decision{Allowed: true, Reason: ReasonAllowed, Snapshot: invalidSnapshot}, want: false},
		{name: "zero snapshot revision", val: Decision{Allowed: true, Reason: ReasonAllowed, Snapshot: zeroRevision}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.val.IsValid(); got != tt.want {
				t.Fatalf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}
