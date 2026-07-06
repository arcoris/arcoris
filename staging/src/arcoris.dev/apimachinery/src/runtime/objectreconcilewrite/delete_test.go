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

package objectreconcilewrite

import "testing"

func TestDeleteBuildsRequestFromCurrent(t *testing.T) {
	key := testKey("task-1")
	current := newCurrent(t, key, 7, "desired")

	req, err := current.Delete()

	requireNoError(t, err)
	if req.Resource != key.Resource {
		t.Fatalf("resource = %#v; want %#v", req.Resource, key.Resource)
	}
	if req.Object != key.Object {
		t.Fatalf("object = %#v; want %#v", req.Object, key.Object)
	}
	if req.Expected != 7 {
		t.Fatalf("expected = %s; want 7", req.Expected)
	}
}

func TestDeleteRejectsZeroCurrent(t *testing.T) {
	var current Current

	_, err := current.Delete()

	requireErrorIs(t, err, ErrInvalidCurrent)
}
