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

package codecselection

import (
	"errors"
	"strings"
	"testing"
)

func TestErrorAtBuildsSelectionError(t *testing.T) {
	err := errorAt(
		"codecselection.decodeBindings[0]",
		ErrInvalidBinding,
		ErrorReasonInvalidBinding,
		"binding is invalid",
	)

	requireErrorIs(t, err, ErrInvalidBinding)
	requireSelectionError(t, err, "codecselection.decodeBindings[0]", ErrorReasonInvalidBinding)
}

func TestErrorfAtBuildsFormattedDetail(t *testing.T) {
	err := errorfAt(
		"codecselection.decodeBindings[0].entryID",
		ErrUnknownEntryID,
		ErrorReasonUnknownEntryID,
		"entry %q is missing",
		"json.public",
	)

	requireSelectionDetailContains(t, err, `entry "json.public" is missing`)
}

func TestWrapAtPreservesCause(t *testing.T) {
	cause := errors.New("cause")
	err := wrapAt(
		"codecselection.decodeBindings[0]",
		ErrInvalidBinding,
		ErrorReasonInvalidBinding,
		"binding is invalid",
		cause,
	)

	requireErrorIs(t, err, ErrInvalidBinding)
	requireErrorIs(t, err, cause)
	requireSelectionDetailContains(t, err, "binding is invalid")
}

func TestWrapAtWithoutCause(t *testing.T) {
	err := wrapAt(
		"codecselection.decodeBindings[0]",
		ErrInvalidBinding,
		ErrorReasonInvalidBinding,
		"binding is invalid",
		nil,
	)

	if strings.Contains(err.Error(), "<nil>") {
		t.Fatalf("error = %q; want no nil cause text", err.Error())
	}
}
