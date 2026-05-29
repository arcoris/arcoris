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

package resource

import (
	"errors"
	"testing"
)

func TestErrorHelpersSetSentinelsAndReasons(t *testing.T) {
	cause := errors.New("cause")

	cases := []struct {
		name   string
		err    error
		target error
		path   string
		reason ErrorReason
	}{
		{
			name:   "definition",
			err:    definitionError("definition", ErrorReasonMissingVersion, "missing"),
			target: ErrInvalidDefinition,
			path:   "definition",
			reason: ErrorReasonMissingVersion,
		},
		{
			name:   "version",
			err:    versionError("version", ErrorReasonMissingDesired, "missing"),
			target: ErrInvalidVersion,
			path:   "version",
			reason: ErrorReasonMissingDesired,
		},
		{
			name:   "scope",
			err:    scopeError(ErrorReasonInvalidScope, detailScopeSupported),
			target: ErrInvalidScope,
			path:   pathScope,
			reason: ErrorReasonInvalidScope,
		},
		{
			name:   "json",
			err:    invalidJSON("scope", detailJSONMustBeString, cause),
			target: ErrInvalidJSON,
			path:   pathScope,
			reason: ErrorReasonInvalidJSON,
		},
		{
			name:   "nil receiver",
			err:    nilReceiver("scope"),
			target: ErrNilReceiver,
			path:   pathScope,
			reason: ErrorReasonNilReceiver,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			requireResourceError(t, tc.err, tc.target, tc.path, tc.reason)
		})
	}
}

func TestNestedErrorHelpersPreserveCauses(t *testing.T) {
	cause := errors.New("cause")

	err := nestedVersionError(
		"definition.versions[v1].desired",
		ErrorReasonInvalidDesired,
		"invalid",
		cause,
	)
	requireErrorIs(t, err, ErrInvalidVersion)
	requireErrorIs(t, err, cause)
}
