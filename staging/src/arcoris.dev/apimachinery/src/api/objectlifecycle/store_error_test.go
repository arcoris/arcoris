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
	"testing"

	"arcoris.dev/apimachinery/api/objectstore"
)

func TestMapStoreErrorPreservesObjectstoreSentinels(t *testing.T) {
	tests := []struct {
		name string
		in   error
		out  error
	}{
		{name: "not found", in: objectstore.ErrNotFound, out: ErrNotFound},
		{name: "already exists", in: objectstore.ErrAlreadyExists, out: ErrAlreadyExists},
		{name: "stale revision", in: objectstore.ErrStaleRevision, out: ErrStaleRevision},
		{name: "conflict", in: objectstore.ErrConflict, out: ErrConflict},
		{name: "other", in: errors.New("other"), out: ErrStoreFailed},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := mapStoreError(OperationApply, objectstore.Key{}, tt.in)

			requireErrorIs(t, err, tt.out)
			requireErrorIs(t, err, tt.in)
		})
	}
}

func TestCreatePreservesObjectstoreAlreadyExistsSentinel(t *testing.T) {
	executor := testExecutor(t)
	createObject(t, executor, 1, "api:v1", owner("creator"))

	_, err := executor.Create(
		context.Background(),
		CreateRequest{Object: testObject(1, "api:v1"), Owner: owner("creator")},
	)

	requireErrorIs(t, err, objectstore.ErrAlreadyExists)
}
