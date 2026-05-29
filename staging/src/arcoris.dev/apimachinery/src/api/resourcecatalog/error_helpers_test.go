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

package resourcecatalog

import (
	"errors"
	"testing"
)

func TestErrorHelpers(t *testing.T) {
	cause := errors.New("cause")

	cases := []struct {
		name   string
		err    error
		target error
		path   string
		reason ErrorReason
	}{
		{
			name:   "catalog",
			err:    catalogError("definitions[0]", ErrDuplicateDefinition, ErrorReasonDuplicateResource, "duplicate"),
			target: ErrDuplicateDefinition,
			path:   "definitions[0]",
			reason: ErrorReasonDuplicateResource,
		},
		{
			name:   "nested",
			err:    nestedCatalogError("definitions[0]", ErrorReasonInvalidDefinition, "invalid", cause),
			target: ErrInvalidCatalog,
			path:   "definitions[0]",
			reason: ErrorReasonInvalidDefinition,
		},
		{
			name:   "nil",
			err:    nilCatalogError(),
			target: ErrNilCatalog,
			path:   "catalog",
			reason: ErrorReasonNilCatalog,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			requireCatalogError(t, tc.err, tc.target, tc.path, tc.reason)
		})
	}

	requireErrorIs(t, cases[1].err, cause)
}
