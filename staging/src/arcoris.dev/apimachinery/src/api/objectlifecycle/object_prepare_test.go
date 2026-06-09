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
	"testing"

	"arcoris.dev/apimachinery/api/value"
)

func TestPrepareObjectRequestResolvesKeysAndValidates(t *testing.T) {
	executor := testExecutor(t)

	prepared, err := executor.prepareObjectRequest(OperationCreate, testObject(1, "api:v1"))
	requireNoError(t, err)

	if prepared.resolved.gvr != testGVR() {
		t.Fatalf("gvr = %#v; want %#v", prepared.resolved.gvr, testGVR())
	}
	if prepared.key.Object != testName(1) {
		t.Fatalf("key object = %#v; want %#v", prepared.key.Object, testName(1))
	}
}

func TestPrepareObjectRequestRejectsInvalidObjectIdentity(t *testing.T) {
	executor := testExecutor(t)
	obj := testObject(1, "api:v1")
	obj.ObjectMeta.Name = ""

	_, err := executor.prepareObjectRequest(OperationCreate, obj)

	requireLifecycleError(t, err, ErrInvalidRequest, ErrorReasonInvalidRequest)
}

func TestPrepareObjectRequestReturnsValidationFailed(t *testing.T) {
	executor := testExecutor(t)
	obj := testObject(1, "api:v1")
	obj.Desired = value.StringValue("not-object")

	_, err := executor.prepareObjectRequest(OperationCreate, obj)

	requireLifecycleError(t, err, ErrValidationFailed, ErrorReasonValidationFailed)
}
