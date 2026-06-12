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

	"arcoris.dev/apimachinery/api/valuevalidation"
)

func TestValidateObjectRequiresObservedValidator(t *testing.T) {
	executor, err := NewExecutor(
		WithStore(testStore(t)),
		WithResourceResolver(testCatalog(t, testDefinition(resourceObserved()))),
		WithDesiredValidator(valuevalidation.SurfaceValidator{}),
	)
	requireNoError(t, err)

	_, err = executor.Create(
		context.Background(),
		CreateRequest{Object: testObservedObject(1, "api:v1", "true"), Owner: owner("creator")},
	)

	requireLifecycleError(t, err, ErrValidationFailed, ErrorReasonValidationFailed)
}

func TestValidateObjectUsesObservedValidator(t *testing.T) {
	executor := testExecutor(
		t,
		WithResourceResolver(testCatalog(t, testDefinition(resourceObserved()))),
	)

	result, err := executor.Create(
		context.Background(),
		CreateRequest{Object: testObservedObject(1, "api:v1", "true"), Owner: owner("creator")},
	)
	requireNoError(t, err)

	if result.State.Object.Observed == nil {
		t.Fatalf("Observed missing")
	}
	requireObservedReady(t, result.State, "true")
}
