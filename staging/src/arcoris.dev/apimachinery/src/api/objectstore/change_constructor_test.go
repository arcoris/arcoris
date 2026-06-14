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

	"arcoris.dev/apimachinery/api/value"
)

func TestNewCreatedChange(t *testing.T) {
	after := committedStateAt(2, "created")

	change, err := NewCreatedChange(validKey(), after)
	requireNoError(t, err)
	after.Object.Desired = value.StringValue("mutated")

	if change.Kind != ChangeCreated || change.Revision != 2 {
		t.Fatalf("change = %#v; want created revision 2", change)
	}
	requireChangeDesired(t, change.After, "created")
	if !isZeroState(change.Before) {
		t.Fatalf("created Before is non-zero")
	}
}

func TestNewUpdatedChange(t *testing.T) {
	before := committedStateAt(1, "before")
	after := committedStateAt(2, "after")

	change, err := NewUpdatedChange(validKey(), before, after)
	requireNoError(t, err)
	before.Object.Desired = value.StringValue("before-mutated")
	after.Object.Desired = value.StringValue("after-mutated")

	if change.Kind != ChangeUpdated || change.Revision != 2 {
		t.Fatalf("change = %#v; want updated revision 2", change)
	}
	requireChangeDesired(t, change.Before, "before")
	requireChangeDesired(t, change.After, "after")
}

func TestNewDeletedChange(t *testing.T) {
	before := committedStateAt(2, "deleted")

	change, err := NewDeletedChange(validKey(), before, 3)
	requireNoError(t, err)
	before.Object.Desired = value.StringValue("mutated")

	if change.Kind != ChangeDeleted || change.Revision != 3 {
		t.Fatalf("change = %#v; want deleted revision 3", change)
	}
	requireChangeDesired(t, change.Before, "deleted")
	if !isZeroState(change.After) {
		t.Fatalf("deleted After is non-zero")
	}
}

func TestChangeConstructorsRejectInvalidInput(t *testing.T) {
	tests := []struct {
		name   string
		build  func() (Change, error)
		target error
	}{
		{
			name: "created invalid key",
			build: func() (Change, error) {
				return NewCreatedChange(Key{}, committedStateAt(1, "after"))
			},
			target: ErrInvalidKey,
		},
		{
			name: "created zero revision",
			build: func() (Change, error) {
				return NewCreatedChange(validKey(), validState())
			},
			target: ErrInvalidRevision,
		},
		{
			name: "updated zero before",
			build: func() (Change, error) {
				return NewUpdatedChange(validKey(), State{}, committedStateAt(2, "after"))
			},
			target: ErrInvalidRevision,
		},
		{
			name: "updated zero after",
			build: func() (Change, error) {
				return NewUpdatedChange(validKey(), committedStateAt(1, "before"), State{})
			},
			target: ErrInvalidRevision,
		},
		{
			name: "deleted zero revision",
			build: func() (Change, error) {
				return NewDeletedChange(validKey(), committedStateAt(1, "before"), 0)
			},
			target: ErrInvalidRevision,
		},
		{
			name: "deleted zero before",
			build: func() (Change, error) {
				return NewDeletedChange(validKey(), State{}, 2)
			},
			target: ErrInvalidRevision,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.build()
			requireErrorIs(t, err, ErrInvalidChange)
			if !errors.Is(err, tt.target) {
				t.Fatalf("errors.Is(%v, %v) = false", err, tt.target)
			}
		})
	}
}

func TestMustChangeConstructorsPanicOnInvalidInput(t *testing.T) {
	tests := []struct {
		name string
		call func()
	}{
		{name: "created", call: func() { _ = MustCreatedChange(Key{}, committedStateAt(1, "after")) }},
		{name: "updated", call: func() { _ = MustUpdatedChange(validKey(), State{}, committedStateAt(2, "after")) }},
		{name: "deleted", call: func() { _ = MustDeletedChange(validKey(), State{}, 2) }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Fatalf("constructor did not panic")
				}
			}()

			tt.call()
		})
	}
}
