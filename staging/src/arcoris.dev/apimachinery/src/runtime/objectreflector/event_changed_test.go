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

package objectreflector

import (
	"context"
	"errors"
	"testing"

	"arcoris.dev/apimachinery/api/objectstore"
)

func TestProcessChangedAppliesCreatedUpdatedDeleted(t *testing.T) {
	tests := []struct {
		name   string
		change objectstore.Change
	}{
		{name: "created", change: createdChange(t, testKey("system", 1), 2)},
		{name: "updated", change: updatedChange(t, testKey("system", 1), 2, 3)},
		{name: "deleted", change: deletedChange(t, testKey("system", 1), 2, 3)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sink := newRecordingSink(1)
			reflector := newTestReflector(t, &fakeListerWatcher{}, sink)
			reflector.lastApplied = tt.change.Revision - 1

			err := reflector.processEvent(context.Background(), changedEvent(t, tt.change))
			requireNoError(t, err)

			if sink.changeCount() != 1 {
				t.Fatalf("change count = %d; want 1", sink.changeCount())
			}
			if reflector.lastApplied != tt.change.Revision {
				t.Fatalf("lastApplied = %s; want %s", reflector.lastApplied, tt.change.Revision)
			}
		})
	}
}

func TestProcessChangedDeliversDetachedChange(t *testing.T) {
	key := testKey("system", 1)
	change := createdChange(t, key, 2)
	event := changedEvent(t, change)
	sink := newRecordingSink(1)
	reflector := newTestReflector(t, &fakeListerWatcher{}, sink)
	reflector.lastApplied = 1

	requireNoError(t, reflector.processEvent(context.Background(), event))
	recorded := sink.recordedChanges()
	recorded[0].After.Revision = 99

	if event.Change.After.Revision != change.Revision {
		t.Fatalf("event payload was mutated through sink copy")
	}
}

func TestProcessChangedDoesNotAdvanceRevisionWhenApplyFails(t *testing.T) {
	errApply := errors.New("apply failed")
	sink := newRecordingSink(1)
	sink.applyErr = errApply
	reflector := newTestReflector(t, &fakeListerWatcher{}, sink)
	reflector.lastApplied = 1

	err := reflector.processEvent(context.Background(), changedEvent(t, createdChange(t, testKey("system", 1), 2)))
	requireErrorIs(t, err, errApply)
	if reflector.lastApplied != 1 {
		t.Fatalf("lastApplied = %s; want 1", reflector.lastApplied)
	}
}

func TestProcessChangedRejectsOutsideCollection(t *testing.T) {
	reflector := newTestReflector(t, &fakeListerWatcher{}, newRecordingSink(1))
	reflector.lastApplied = 1

	err := reflector.processEvent(context.Background(), changedEvent(t, createdChange(t, otherResourceKey("system", 1), 2)))

	requireErrorIs(t, err, ErrChangeOutsideCollection)
}

func TestProcessChangedRejectsNonMonotonicRevision(t *testing.T) {
	reflector := newTestReflector(t, &fakeListerWatcher{}, newRecordingSink(1))
	reflector.lastApplied = 2

	err := reflector.processEvent(context.Background(), changedEvent(t, createdChange(t, testKey("system", 1), 2)))

	requireErrorIs(t, err, ErrNonMonotonicRevision)
}
