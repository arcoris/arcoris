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

func TestChangeValidate(t *testing.T) {
	tests := []struct {
		name       string
		change     Change
		wantReason ErrorReason
		wantCause  error
	}{
		{
			name: "valid created",
			change: Change{
				Kind:     ChangeCreated,
				Key:      validKey(),
				Revision: 1,
				After:    committedStateAt(1, "after"),
			},
		},
		{
			name: "created with non-zero before",
			change: Change{
				Kind:     ChangeCreated,
				Key:      validKey(),
				Revision: 2,
				Before:   committedStateAt(1, "before"),
				After:    committedStateAt(2, "after"),
			},
			wantReason: ErrorReasonInvalidChange,
			wantCause:  ErrInvalidChange,
		},
		{
			name: "created with zero after",
			change: Change{
				Kind:     ChangeCreated,
				Key:      validKey(),
				Revision: 1,
			},
			wantReason: ErrorReasonInvalidChange,
			wantCause:  ErrInvalidRevision,
		},
		{
			name: "created after revision mismatch",
			change: Change{
				Kind:     ChangeCreated,
				Key:      validKey(),
				Revision: 2,
				After:    committedStateAt(1, "after"),
			},
			wantReason: ErrorReasonInvalidChange,
			wantCause:  ErrInvalidRevision,
		},
		{
			name: "valid updated",
			change: Change{
				Kind:     ChangeUpdated,
				Key:      validKey(),
				Revision: 2,
				Before:   committedStateAt(1, "before"),
				After:    committedStateAt(2, "after"),
			},
		},
		{
			name: "updated with zero before",
			change: Change{
				Kind:     ChangeUpdated,
				Key:      validKey(),
				Revision: 2,
				After:    committedStateAt(2, "after"),
			},
			wantReason: ErrorReasonInvalidChange,
			wantCause:  ErrInvalidRevision,
		},
		{
			name: "updated with zero after",
			change: Change{
				Kind:     ChangeUpdated,
				Key:      validKey(),
				Revision: 2,
				Before:   committedStateAt(1, "before"),
			},
			wantReason: ErrorReasonInvalidChange,
			wantCause:  ErrInvalidRevision,
		},
		{
			name: "updated after revision mismatch",
			change: Change{
				Kind:     ChangeUpdated,
				Key:      validKey(),
				Revision: 3,
				Before:   committedStateAt(1, "before"),
				After:    committedStateAt(2, "after"),
			},
			wantReason: ErrorReasonInvalidChange,
			wantCause:  ErrInvalidRevision,
		},
		{
			name: "updated with non-newer after revision",
			change: Change{
				Kind:     ChangeUpdated,
				Key:      validKey(),
				Revision: 1,
				Before:   committedStateAt(1, "before"),
				After:    committedStateAt(1, "after"),
			},
			wantReason: ErrorReasonInvalidChange,
			wantCause:  ErrInvalidRevision,
		},
		{
			name: "valid deleted",
			change: Change{
				Kind:     ChangeDeleted,
				Key:      validKey(),
				Revision: 2,
				Before:   committedStateAt(1, "before"),
			},
		},
		{
			name: "deleted with zero before",
			change: Change{
				Kind:     ChangeDeleted,
				Key:      validKey(),
				Revision: 2,
			},
			wantReason: ErrorReasonInvalidChange,
			wantCause:  ErrInvalidRevision,
		},
		{
			name: "deleted with non-zero after",
			change: Change{
				Kind:     ChangeDeleted,
				Key:      validKey(),
				Revision: 2,
				Before:   committedStateAt(1, "before"),
				After:    committedStateAt(2, "after"),
			},
			wantReason: ErrorReasonInvalidChange,
			wantCause:  ErrInvalidChange,
		},
		{
			name: "deleted with zero tombstone revision",
			change: Change{
				Kind:   ChangeDeleted,
				Key:    validKey(),
				Before: committedStateAt(1, "before"),
			},
			wantReason: ErrorReasonInvalidChange,
			wantCause:  ErrInvalidRevision,
		},
		{
			name: "deleted with non-newer tombstone revision",
			change: Change{
				Kind:     ChangeDeleted,
				Key:      validKey(),
				Revision: 1,
				Before:   committedStateAt(1, "before"),
			},
			wantReason: ErrorReasonInvalidChange,
			wantCause:  ErrInvalidRevision,
		},
		{
			name: "unknown kind",
			change: Change{
				Kind:     ChangeKind(99),
				Key:      validKey(),
				Revision: 1,
				After:    committedStateAt(1, "after"),
			},
			wantReason: ErrorReasonInvalidChangeKind,
			wantCause:  ErrInvalidChange,
		},
		{
			name:       "zero change",
			change:     Change{},
			wantReason: ErrorReasonInvalidChangeKind,
			wantCause:  ErrInvalidChange,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateChange(tt.change)
			if tt.wantCause == nil {
				requireNoError(t, err)
				return
			}

			requireErrorIs(t, err, ErrInvalidChange)
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
		})
	}
}
