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
	"testing"
)

func TestProcessProgressDoesNotApplyChange(t *testing.T) {
	sink := newRecordingSink(1)
	reflector := newTestReflector(t, &fakeListerWatcher{}, sink)
	reflector.lastApplied = 2
	reflector.lastProgress = 2

	err := reflector.processEvent(context.Background(), progressEvent(t, 3))
	requireNoError(t, err)

	if sink.changeCount() != 0 {
		t.Fatalf("change count = %d; want 0", sink.changeCount())
	}
	if reflector.lastProgress != 3 {
		t.Fatalf("lastProgress = %s; want 3", reflector.lastProgress)
	}
}

func TestProcessProgressRejectsOlderRevision(t *testing.T) {
	reflector := newTestReflector(t, &fakeListerWatcher{}, newRecordingSink(1))
	reflector.lastApplied = 3
	reflector.lastProgress = 3

	err := reflector.processEvent(context.Background(), progressEvent(t, 2))

	requireErrorIs(t, err, ErrNonMonotonicRevision)
}
