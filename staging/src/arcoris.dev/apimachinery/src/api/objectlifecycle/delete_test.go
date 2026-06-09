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
	"testing"
)

func TestDeleteExistingRevisionDeletesState(t *testing.T) {
	executor := testExecutor(t)
	created := createObject(t, executor, 1, "api:v1", owner("creator"))

	result, err := executor.Delete(
		context.Background(),
		DeleteRequest{Resource: testGVR(), Object: testName(1), Expected: created.State.Revision},
	)
	requireNoError(t, err)

	requireEffect(t, result, OperationDelete, EffectDeleted)
	if result.State.Revision != created.State.Revision {
		t.Fatalf("deleted revision = %v; want %v", result.State.Revision, created.State.Revision)
	}

	_, err = executor.Get(context.Background(), GetRequest{Resource: testGVR(), Object: testName(1)})
	requireLifecycleError(t, err, ErrNotFound, ErrorReasonNotFound)
}

func TestDeleteMissingReturnsNotFound(t *testing.T) {
	executor := testExecutor(t)

	_, err := executor.Delete(
		context.Background(),
		DeleteRequest{Resource: testGVR(), Object: testName(1), Expected: 1},
	)

	requireLifecycleError(t, err, ErrNotFound, ErrorReasonNotFound)
}

func TestDeleteStaleRevisionReturnsStaleRevision(t *testing.T) {
	executor := testExecutor(t)
	createObject(t, executor, 1, "api:v1", owner("creator"))

	_, err := executor.Delete(
		context.Background(),
		DeleteRequest{Resource: testGVR(), Object: testName(1), Expected: 99},
	)

	requireLifecycleError(t, err, ErrStaleRevision, ErrorReasonStaleRevision)
}

func TestDeleteZeroExpectedRevisionReturnsInvalidRequest(t *testing.T) {
	executor := testExecutor(t)

	_, err := executor.Delete(
		context.Background(),
		DeleteRequest{Resource: testGVR(), Object: testName(1)},
	)

	requireLifecycleError(t, err, ErrInvalidRequest, ErrorReasonInvalidRequest)
}
