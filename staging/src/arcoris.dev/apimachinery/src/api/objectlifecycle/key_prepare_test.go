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

import "testing"

func TestPrepareKeyRequestResolvesResourceAndKey(t *testing.T) {
	executor := testExecutor(t)

	prepared, err := executor.prepareKeyRequest(OperationGet, testGVR(), testName(1))
	requireNoError(t, err)

	if prepared.resolved.gvr != testGVR() {
		t.Fatalf("gvr = %#v; want %#v", prepared.resolved.gvr, testGVR())
	}
	if prepared.key.Object != testName(1) {
		t.Fatalf("key object = %#v; want %#v", prepared.key.Object, testName(1))
	}
}

func TestPrepareKeyRequestReturnsResourceNotFound(t *testing.T) {
	executor := testExecutor(t)
	gvr := testGVR()
	gvr.Resource = "unknowns"

	_, err := executor.prepareKeyRequest(OperationGet, gvr, testName(1))

	requireLifecycleError(t, err, ErrResourceNotFound, ErrorReasonResourceNotFound)
}

func TestPrepareKeyRequestRejectsInvalidObjectName(t *testing.T) {
	executor := testExecutor(t)
	name := testName(1)
	name.Name = ""

	_, err := executor.prepareKeyRequest(OperationGet, testGVR(), name)

	requireLifecycleError(t, err, ErrInvalidRequest, ErrorReasonInvalidRequest)
}
