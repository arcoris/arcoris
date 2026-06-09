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

package objectlifecycle

import (
	"context"
	"errors"
	"strings"
	"testing"

	"arcoris.dev/apimachinery/api/objectstore"
)

func TestErrorSupportsErrorsIsForLifecycleAndCause(t *testing.T) {
	err := errorFor(OperationApply, ReasonStaleRevision, objectstore.Key{}, ErrStaleRevision, objectstore.ErrStaleRevision)

	requireErrorIs(t, err, ErrStaleRevision)
	requireErrorIs(t, err, objectstore.ErrStaleRevision)
}

func TestErrorStringIncludesOperationReasonAndKey(t *testing.T) {
	key := objectstore.MustKey(testGVR(), testName(1))
	err := errorFor(OperationDelete, ReasonStaleRevision, key, ErrStaleRevision, objectstore.ErrStaleRevision).Error()

	for _, want := range []string{"delete", "stale_revision", key.String()} {
		if !strings.Contains(err, want) {
			t.Fatalf("Error() = %q; want to contain %q", err, want)
		}
	}
}

func TestNilErrorIsSafe(t *testing.T) {
	var err *Error
	if got := err.Error(); got != "<nil>" {
		t.Fatalf("Error() = %q; want <nil>", got)
	}
	if got := err.Unwrap(); got != nil {
		t.Fatalf("Unwrap() = %v; want nil", got)
	}
}

func TestCreatePreservesObjectstoreAlreadyExistsSentinel(t *testing.T) {
	executor := testExecutor(t)
	createObject(t, executor, 1, "api:v1", owner("creator"))

	_, err := executor.Create(
		context.Background(),
		CreateRequest{Object: testObject(1, "api:v1"), Owner: owner("creator")},
	)

	if !errors.Is(err, objectstore.ErrAlreadyExists) {
		t.Fatalf("errors.Is(%v, objectstore.ErrAlreadyExists) = false", err)
	}
}
