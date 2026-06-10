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

package valuemerge

import (
	"errors"
	"testing"
)

func TestErrorIsSentinel(t *testing.T) {
	err := errorAt(root(), ErrUnsupportedMerge, ErrorReasonUnsupportedMerge, "nope")

	requireErrorIs(t, err, ErrUnsupportedMerge)
}

func TestErrorAsValueMergeError(t *testing.T) {
	err := errorAt(root(), ErrUnsupportedMerge, ErrorReasonUnsupportedMerge, "nope")

	var mergeError *Error
	if !errors.As(err, &mergeError) {
		t.Fatalf("errors.As(*Error) = false")
	}
	if mergeError.Reason != ErrorReasonUnsupportedMerge {
		t.Fatalf("reason = %q; want %q", mergeError.Reason, ErrorReasonUnsupportedMerge)
	}
}
