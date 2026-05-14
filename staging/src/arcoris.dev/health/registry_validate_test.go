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

package health

import (
	"context"
	"errors"
	"testing"
)

func TestPrepareChecksCollectsBatchFailures(t *testing.T) {
	t.Parallel()

	var typedNil *typedNilChecker
	invalid := checkerFunc{name: "bad-name", fn: func(context.Context) Result { return Healthy("bad-name") }}
	firstDuplicate := mustCheck(t, "duplicate", Healthy("duplicate"))
	secondDuplicate := mustCheck(t, "duplicate", Healthy("duplicate"))

	prepared, err := prepareChecks(
		TargetReady,
		[]Checker{
			nil,
			typedNil,
			invalid,
			firstDuplicate,
			secondDuplicate,
		},
	)
	if prepared != nil {
		t.Fatalf("prepared = %+v, want nil on validation failure", prepared)
	}
	for _, target := range []error{ErrNilChecker, ErrInvalidCheckName, ErrDuplicateCheck} {
		if !errors.Is(err, target) {
			t.Fatalf("errors.Is(%v, %v) = false, want true", err, target)
		}
	}

	children := joinedErrors(err)
	if len(children) != 4 {
		t.Fatalf("joined validation errors = %d, want 4", len(children))
	}
}

func TestPrepareChecksReturnsPreparedMetadata(t *testing.T) {
	t.Parallel()

	first := mustCheck(t, "first", Healthy("first"))
	second := mustCheck(t, "second", Healthy("second"))

	prepared, err := prepareChecks(TargetReady, []Checker{first, second})
	if err != nil {
		t.Fatalf("prepareChecks() = %v, want nil", err)
	}
	if len(prepared) != 2 {
		t.Fatalf("prepared length = %d, want 2", len(prepared))
	}
	if prepared[0].Index != 0 || prepared[0].Name != "first" || prepared[0].Checker.Name() != first.Name() {
		t.Fatalf("prepared[0] = %+v, want first metadata", prepared[0])
	}
	if prepared[1].Index != 1 || prepared[1].Name != "second" || prepared[1].Checker.Name() != second.Name() {
		t.Fatalf("prepared[1] = %+v, want second metadata", prepared[1])
	}
}

func TestNilCheckerDetectsTypedNil(t *testing.T) {
	t.Parallel()

	var typedNil *typedNilChecker
	if !nilChecker(nil) {
		t.Fatal("nilChecker(nil) = false, want true")
	}
	if !nilChecker(typedNil) {
		t.Fatal("nilChecker(typed nil) = false, want true")
	}
	if nilChecker(mustCheck(t, "storage", Healthy("storage"))) {
		t.Fatal("nilChecker(valid) = true, want false")
	}
}
