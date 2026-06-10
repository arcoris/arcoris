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

package fieldownership

import (
	"errors"
	"testing"

	"arcoris.dev/apimachinery/api/fieldpath"
)

func TestNewEntry(t *testing.T) {
	entry, err := NewEntry(owner("user-cli"), set(imagePath()))

	requireNoError(t, err)
	requireEqual(t, entry.Owner(), owner("user-cli"))
	requireSet(t, entry.Fields(), "$.spec.image")
}

func TestNewEntryRejectsInvalidOwner(t *testing.T) {
	_, err := NewEntry(Owner{}, set(imagePath()))

	requireErrorIs(t, err, ErrInvalidEntry)
}

func TestNewEntryInvalidOwnerIsInvalidEntryAndInvalidOwner(t *testing.T) {
	_, err := NewEntry(Owner{}, set(imagePath()))

	requireErrorIs(t, err, ErrInvalidEntry)
	requireErrorIs(t, err, ErrInvalidOwner)

	var fieldError *Error
	if !errors.As(err, &fieldError) {
		t.Fatalf("errors.As did not find Error")
	}
}

func TestNewEntryAllowsEmptyFields(t *testing.T) {
	entry, err := NewEntry(owner("user-cli"), fieldpath.EmptySet())

	requireNoError(t, err)
	requireEqual(t, entry.IsEmpty(), true)
}

func TestMustEntryPanicsOnInvalidOwner(t *testing.T) {
	requirePanic(t, func() {
		MustEntry(Owner{}, set(imagePath()))
	})
}
