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
	"errors"
	"testing"
)

func TestZeroResultIsSuccess(t *testing.T) {
	var result Result

	if result.Failed() {
		t.Fatalf("zero Result Failed() = true; want false")
	}
}

func TestResultWithErrorIsFailure(t *testing.T) {
	err := errors.New("failed")
	result := Result{Err: err}

	if !result.Failed() {
		t.Fatalf("Result{Err} Failed() = false; want true")
	}
	if result.Err != err {
		t.Fatalf("Err = %v; want %v", result.Err, err)
	}
}

func TestResultHelpers(t *testing.T) {
	if Success().Failed() {
		t.Fatalf("Success().Failed() = true; want false")
	}
	if Failure(nil).Failed() {
		t.Fatalf("Failure(nil).Failed() = true; want false")
	}

	err := errors.New("failed")
	failure := Failure(err)
	if !failure.Failed() || failure.Err != err {
		t.Fatalf("Failure(err) = %#v; want err", failure)
	}
}
