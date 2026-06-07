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

package health

import (
	"errors"
	"testing"
)

func TestCheckSetConstructsOrderedImmutableSet(t *testing.T) {
	t.Parallel()

	storage := mustCheck(t, "storage", Healthy("storage"))
	queue := mustCheck(t, "queue", Healthy("queue"))

	set, err := NewCheckSet(TargetReady, storage, queue)
	if err != nil {
		t.Fatalf("NewCheckSet() = %v, want nil", err)
	}
	if !set.IsValid() || set.Target() != TargetReady || set.Len() != 2 || set.Empty() {
		t.Fatalf("set invariants = %+v", set)
	}
	if !set.Has("storage") || set.Has("missing") || set.Has("bad-name") {
		t.Fatalf("Has() mismatch for set")
	}

	got := set.Checks()
	got[0] = queue
	if set.Checks()[0].Name() != "storage" {
		t.Fatal("Checks() exposed mutable internal storage")
	}

	var names []string
	set.Range(func(checker Checker) bool {
		names = append(names, checker.Name())
		return true
	})
	if len(names) != 2 || names[0] != "storage" || names[1] != "queue" {
		t.Fatalf("Range order = %v, want storage queue", names)
	}
}

func TestCheckSetRejectsInvalidInputs(t *testing.T) {
	t.Parallel()

	valid := mustCheck(t, "storage", Healthy("storage"))

	tests := []struct {
		name string
		call func() error
		want error
	}{
		{
			name: "unknown target",
			call: func() error {
				_, err := NewCheckSet(TargetUnknown, valid)
				return err
			},
			want: ErrInvalidTarget,
		},
		{
			name: "nil checker",
			call: func() error {
				_, err := NewCheckSet(TargetReady, nil)
				return err
			},
			want: ErrNilChecker,
		},
		{
			name: "typed nil checker",
			call: func() error {
				var checker *typedNilChecker
				_, err := NewCheckSet(TargetReady, checker)
				return err
			},
			want: ErrNilChecker,
		},
		{
			name: "invalid name",
			call: func() error {
				_, err := NewCheckSet(TargetReady, checkerFunc{name: "bad-name"})
				return err
			},
			want: ErrInvalidCheckName,
		},
		{
			name: "duplicate name",
			call: func() error {
				_, err := NewCheckSet(TargetReady, valid, valid)
				return err
			},
			want: ErrDuplicateCheckName,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if err := tc.call(); !errors.Is(err, tc.want) {
				t.Fatalf("NewCheckSet() = %v, want %v", err, tc.want)
			}
		})
	}
}

func TestCheckSetDuplicateCarriesIndexes(t *testing.T) {
	t.Parallel()

	checker := mustCheck(t, "dup", Healthy("dup"))
	_, err := NewCheckSet(TargetReady, checker, checker)

	var duplicate DuplicateCheckNameError
	if !errors.As(err, &duplicate) {
		t.Fatalf("NewCheckSet() = %v, want DuplicateCheckNameError", err)
	}
	if duplicate.Name != "dup" || duplicate.Index != 1 || duplicate.PreviousIndex != 0 {
		t.Fatalf("duplicate = %+v, want dup index 1 previous 0", duplicate)
	}
}

func TestCheckSetEmptyConcreteTargetIsValid(t *testing.T) {
	t.Parallel()

	set, err := NewCheckSet(TargetLive)
	if err != nil {
		t.Fatalf("NewCheckSet(empty) = %v, want nil", err)
	}
	if !set.IsValid() || !set.Empty() || set.Target() != TargetLive {
		t.Fatalf("empty set = %+v, want valid live empty set", set)
	}
}
