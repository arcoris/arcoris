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

package objectreflector

import "testing"

func TestDefaultOptionsRequestsProgress(t *testing.T) {
	options := DefaultOptions()

	if !options.RequestProgress {
		t.Fatalf("RequestProgress = false; want true")
	}
}

func TestWithRequestProgress(t *testing.T) {
	options := DefaultOptions()
	err := WithRequestProgress(false)(&options)
	requireNoError(t, err)

	if options.RequestProgress {
		t.Fatalf("RequestProgress = true; want false")
	}
}

func TestApplyOptionsRejectsNilOption(t *testing.T) {
	_, err := applyOptions([]Option{nil})

	requireErrorIs(t, err, ErrInvalidOption)
}
