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
	"testing"

	"arcoris.dev/apimachinery/api/types"
)

func TestValidateSurfaceRejectsInvalidAndNonObjectSurfaces(t *testing.T) {
	cases := []struct {
		name   string
		typ    types.Type
		path   string
		target error
		reason ErrorReason
	}{
		{
			name:   "invalid desired type",
			typ:    types.Type{},
			path:   "definition.versions[v1].desired",
			target: ErrInvalidVersion,
			reason: ErrorReasonInvalidDesired,
		},
		{
			name:   "desired scalar",
			typ:    stringType(),
			path:   "definition.versions[v1].desired",
			target: ErrInvalidVersion,
			reason: ErrorReasonDesiredNotObject,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateSurface(
				tc.typ,
				nil,
				tc.path,
				ErrorReasonInvalidDesired,
				ErrorReasonDesiredNotObject,
				detailDesiredObjectLikeTemplate,
			)
			requireResourceError(t, err, tc.target, tc.path, tc.reason)
		})
	}
}
