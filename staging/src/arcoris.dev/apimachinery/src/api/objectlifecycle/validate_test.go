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

	"arcoris.dev/apimachinery/api/fieldownership"
)

func TestCreateRejectsNilContext(t *testing.T) {
	executor := testExecutor(t)

	_, err := executor.Create(nil, CreateRequest{Object: testObject(1, "api:v1"), Owner: owner("creator")})

	requireLifecycleError(t, err, ErrInvalidRequest, ErrorReasonInvalidContext)
	requireErrorIs(t, err, ErrNilContext)
}

func TestCreateRejectsInvalidOwner(t *testing.T) {
	executor := testExecutor(t)

	_, err := executor.Create(
		context.Background(),
		CreateRequest{Object: testObject(1, "api:v1")},
	)

	requireLifecycleError(t, err, ErrInvalidRequest, ErrorReasonInvalidOwner)
	requireErrorIs(t, err, fieldownership.ErrInvalidOwner)
}
