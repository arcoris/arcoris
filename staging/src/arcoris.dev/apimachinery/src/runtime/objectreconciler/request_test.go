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

package objectreconciler

import (
	"testing"

	"arcoris.dev/apimachinery/api/objectstore"
)

func TestRequestValidateAcceptsValidKey(t *testing.T) {
	request := testRequest(1)

	err := request.Validate()

	requireNoError(t, err)
	if !request.IsValid() {
		t.Fatalf("IsValid() = false; want true")
	}
}

func TestRequestValidateRejectsInvalidKey(t *testing.T) {
	request := Request{}

	err := request.Validate()

	requireErrorIs(t, err, ErrInvalidRequest)
	requireErrorIs(t, err, objectstore.ErrInvalidKey)
	if request.IsValid() {
		t.Fatalf("IsValid() = true; want false")
	}
}
