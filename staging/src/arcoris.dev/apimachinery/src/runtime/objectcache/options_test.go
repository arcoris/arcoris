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

package objectcache

import (
	"errors"
	"testing"
)

func TestApplyOptionsAcceptsEmptyOptions(t *testing.T) {
	options, err := applyOptions(nil)
	requireNoError(t, err)
	if options.History.RetainedVersionsPerObject != 0 {
		t.Fatalf("RetainedVersionsPerObject = %d; want 0", options.History.RetainedVersionsPerObject)
	}
}

func TestApplyOptionsRejectsNilOption(t *testing.T) {
	_, err := applyOptions([]Option{nil})
	requireErrorIs(t, err, ErrInvalidOption)
}

func TestApplyOptionsPreservesOptionError(t *testing.T) {
	cause := errors.New("option failed")
	_, err := applyOptions([]Option{func(*Options) error { return cause }})

	requireErrorIs(t, err, cause)
}

func TestDefaultOptionsDisablesHistory(t *testing.T) {
	options := DefaultOptions()
	if options.History.RetainedVersionsPerObject != 0 {
		t.Fatalf("RetainedVersionsPerObject = %d; want 0", options.History.RetainedVersionsPerObject)
	}
}

func TestWithHistoryValidatesRetention(t *testing.T) {
	for _, retained := range []int{0, 1, 2, 8} {
		options, err := applyOptions([]Option{WithHistory(HistoryPolicy{RetainedVersionsPerObject: retained})})
		requireNoError(t, err)
		if options.History.RetainedVersionsPerObject != retained {
			t.Fatalf("RetainedVersionsPerObject = %d; want %d", options.History.RetainedVersionsPerObject, retained)
		}
	}

	_, err := applyOptions([]Option{WithHistory(HistoryPolicy{RetainedVersionsPerObject: -1})})
	requireErrorIs(t, err, ErrInvalidOption)
}
