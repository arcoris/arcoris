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

package healthregistry

import (
	"errors"
	"testing"

	"arcoris.dev/health"
)

func TestRegistryErrorsClassifyRootAndRegistrySentinels(t *testing.T) {
	t.Parallel()

	nilErr := NilCheckerError{Target: health.TargetReady, Index: 1}
	if !errors.Is(nilErr, health.ErrNilChecker) {
		t.Fatal("NilCheckerError should match health.ErrNilChecker")
	}

	nameErr := InvalidCheckNameError{
		Target: health.TargetReady,
		Index:  1,
		Name:   "bad-name",
		Err:    health.ErrInvalidCheckName,
	}
	if !errors.Is(nameErr, health.ErrInvalidCheckName) {
		t.Fatal("InvalidCheckNameError should unwrap health.ErrInvalidCheckName")
	}

	dupErr := DuplicateCheckError{Target: health.TargetReady, Name: "storage", Index: 2, PreviousIndex: 0}
	if !errors.Is(dupErr, ErrDuplicateCheck) {
		t.Fatal("DuplicateCheckError should match ErrDuplicateCheck")
	}
}
