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

package capacity_test

import (
	"errors"
	"testing"

	"arcoris.dev/capacity"
)

func TestNewVectorCanonicalizesAndCopies(t *testing.T) {
	t.Parallel()

	input := []capacity.Entry{
		entry("memory_bytes", 8),
		entry("worker_slots", 2),
	}
	v := vector(t, input...)
	input[0] = entry("queue_slots", 1)

	requireVector(t, v, entry("memory_bytes", 8), entry("worker_slots", 2))
	if !v.IsValid() || v.IsZero() || v.Len() != 2 {
		t.Fatalf("vector validity/length mismatch: valid=%v zero=%v len=%d", v.IsValid(), v.IsZero(), v.Len())
	}
}

func TestNewVectorRejectsInvalidInput(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		entries  []capacity.Entry
		sentinel error
	}{
		{
			name:     "invalid resource",
			entries:  []capacity.Entry{{Resource: "bad-name", Amount: 1}},
			sentinel: capacity.ErrInvalidResource,
		},
		{
			name:     "zero amount",
			entries:  []capacity.Entry{{Resource: "worker_slots"}},
			sentinel: capacity.ErrZeroAmount,
		},
		{
			name: "duplicate resource",
			entries: []capacity.Entry{
				entry("worker_slots", 1),
				entry("worker_slots", 2),
			},
			sentinel: capacity.ErrDuplicateResource,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := capacity.NewVector(tt.entries...)
			if !errors.Is(err, tt.sentinel) {
				t.Fatalf("NewVector() error = %v, want errors.Is(..., %v)", err, tt.sentinel)
			}
			var capacityErr *capacity.Error
			if !errors.As(err, &capacityErr) {
				t.Fatalf("NewVector() error = %T, want *capacity.Error", err)
			}
			if capacityErr.Path == "" {
				t.Fatal("diagnostic path is empty")
			}
		})
	}
}
