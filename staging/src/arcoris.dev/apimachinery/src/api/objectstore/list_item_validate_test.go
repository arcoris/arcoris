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

package objectstore

import (
	"errors"
	"testing"
)

func TestValidateListItem(t *testing.T) {
	tests := []struct {
		name       string
		item       ListItem
		wantReason ErrorReason
		wantCause  error
	}{
		{
			name: "valid",
			item: ListItem{Key: validKey(), State: validCommittedState()},
		},
		{
			name:       "invalid key",
			item:       ListItem{State: validCommittedState()},
			wantReason: ErrorReasonInvalidListItem,
			wantCause:  ErrInvalidKey,
		},
		{
			name:       "invalid state",
			item:       ListItem{Key: validKey(), State: State{}},
			wantReason: ErrorReasonInvalidListItem,
			wantCause:  ErrInvalidRevision,
		},
		{
			name: "state identity mismatch",
			item: ListItem{
				Key:   validKey(),
				State: committedStateForObject(1, "other", "main", "desired"),
			},
			wantReason: ErrorReasonInvalidListItem,
			wantCause:  ErrInvalidKey,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateListItem(tt.item)
			if tt.wantCause == nil {
				requireNoError(t, err)
				if !tt.item.IsValid() {
					t.Fatalf("IsValid() = false; want true")
				}
				return
			}

			requireErrorIs(t, err, ErrInvalidListResult)
			if !errors.Is(err, tt.wantCause) {
				t.Fatalf("errors.Is(%v, %v) = false", err, tt.wantCause)
			}
			var storeErr *Error
			if !errors.As(err, &storeErr) {
				t.Fatalf("error = %T; want *Error", err)
			}
			if storeErr.Reason != tt.wantReason {
				t.Fatalf("reason = %q; want %q", storeErr.Reason, tt.wantReason)
			}
			if tt.item.IsValid() {
				t.Fatalf("IsValid() = true; want false")
			}
		})
	}
}
