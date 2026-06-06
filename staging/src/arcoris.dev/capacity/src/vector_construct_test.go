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

func TestNewVectorSortsAndCopiesEntries(t *testing.T) {
	input := []capacity.Entry{
		entry("worker_slots", 2),
		entry("memory_bytes", 8),
	}

	vector := vector(t, input...)
	input[0] = entry("queue_slots", 1)

	requireVector(t, vector, entry("memory_bytes", 8), entry("worker_slots", 2))
}

func TestNewVectorRejectsInvalidEntries(t *testing.T) {
	tests := []struct {
		name    string
		entries []capacity.Entry
		err     error
	}{
		{name: "invalid resource", entries: []capacity.Entry{{Resource: "bad-name", Amount: 1}}, err: capacity.ErrInvalidResource},
		{name: "zero amount", entries: []capacity.Entry{{Resource: "worker_slots"}}, err: capacity.ErrZeroAmount},
		{name: "duplicate resource", entries: []capacity.Entry{entry("worker_slots", 1), entry("worker_slots", 2)}, err: capacity.ErrDuplicateResource},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := capacity.NewVector(tt.entries...)
			if !errors.Is(err, tt.err) {
				t.Fatalf("NewVector() error = %v, want %v", err, tt.err)
			}
		})
	}
}
