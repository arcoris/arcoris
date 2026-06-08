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

package memory

import (
	"errors"
	"testing"

	"arcoris.dev/apimachinery/api/objectstore"
)

func TestObjectstoreErrorBuildsStructuredError(t *testing.T) {
	key := testKey(1)
	err := objectstoreError(objectstore.ErrorReasonConflict, key, 1, 2, objectstore.ErrConflict)

	var storeErr *objectstore.Error
	if !errors.As(err, &storeErr) {
		t.Fatalf("errors.As failed")
	}
	if storeErr.Reason != objectstore.ErrorReasonConflict || !storeErr.Key.Equal(key) {
		t.Fatalf("structured error = %#v", storeErr)
	}
	if storeErr.Expected != 1 || storeErr.Actual != 2 {
		t.Fatalf("revisions = expected %v actual %v; want 1 and 2", storeErr.Expected, storeErr.Actual)
	}
	requireErrorIs(t, err, objectstore.ErrConflict)
}
