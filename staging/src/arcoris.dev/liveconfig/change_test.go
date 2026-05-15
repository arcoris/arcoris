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

package liveconfig

import (
	"testing"

	"arcoris.dev/snapshot"
)

func TestChangeStateMethods(t *testing.T) {
	tests := []struct {
		name      string
		change    Change[string]
		isChanged bool
		isNoop    bool
		accepted  bool
		rejected  bool
	}{
		{
			name:      "published",
			change:    Change[string]{Changed: true, Reason: ChangeReasonPublished},
			isChanged: true,
			isNoop:    false,
			accepted:  true,
		},
		{
			name:      "equal",
			change:    Change[string]{Changed: false, Reason: ChangeReasonEqual},
			isChanged: false,
			isNoop:    true,
			accepted:  true,
		},
		{
			name:      "normalize failed",
			change:    Change[string]{Changed: false, Reason: ChangeReasonNormalizeFailed},
			isChanged: false,
			isNoop:    true,
			rejected:  true,
		},
		{
			name:      "validate failed",
			change:    Change[string]{Changed: false, Reason: ChangeReasonValidateFailed},
			isChanged: false,
			isNoop:    true,
			rejected:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.change.IsChanged(); got != tt.isChanged {
				t.Fatalf("IsChanged() = %v, want %v", got, tt.isChanged)
			}
			if got := tt.change.IsNoop(); got != tt.isNoop {
				t.Fatalf("IsNoop() = %v, want %v", got, tt.isNoop)
			}
			if got := tt.change.Accepted(); got != tt.accepted {
				t.Fatalf("Accepted() = %v, want %v", got, tt.accepted)
			}
			if got := tt.change.Rejected(); got != tt.rejected {
				t.Fatalf("Rejected() = %v, want %v", got, tt.rejected)
			}
		})
	}
}

func TestChangeCarriesPreviousAndCurrentSnapshots(t *testing.T) {
	change := Change[string]{
		Previous: snapshot.Snapshot[string]{Revision: 1, Value: "prev"},
		Current:  snapshot.Snapshot[string]{Revision: 2, Value: "cur"},
		Changed:  true,
		Reason:   ChangeReasonPublished,
	}

	if got, want := change.Previous.Value, "prev"; got != want {
		t.Fatalf("Previous.Value = %q, want %q", got, want)
	}
	if got, want := change.Current.Value, "cur"; got != want {
		t.Fatalf("Current.Value = %q, want %q", got, want)
	}
}

func TestChangeRevisionRelationship(t *testing.T) {
	tests := []struct {
		name   string
		change Change[string]
	}{
		{
			name: "changed",
			change: Change[string]{
				Previous: snapshot.Snapshot[string]{Revision: 1, Value: "prev"},
				Current:  snapshot.Snapshot[string]{Revision: 2, Value: "cur"},
				Changed:  true,
				Reason:   ChangeReasonPublished,
			},
		},
		{
			name: "equal",
			change: Change[string]{
				Previous: snapshot.Snapshot[string]{Revision: 1, Value: "prev"},
				Current:  snapshot.Snapshot[string]{Revision: 1, Value: "prev"},
				Changed:  false,
				Reason:   ChangeReasonEqual,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sameRevision := tt.change.Current.Revision == tt.change.Previous.Revision
			if tt.change.Changed == sameRevision {
				t.Fatalf("Changed = %v with sameRevision = %v", tt.change.Changed, sameRevision)
			}
		})
	}
}
