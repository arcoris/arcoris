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

func TestValidateSurfaceLocalRejectsInvalidAndNonObjectSurfaces(t *testing.T) {
	cases := []struct {
		name   string
		desc   types.Descriptor
		path   string
		target error
		reason ErrorReason
	}{
		{
			name:   "invalid desired type",
			desc:   types.Descriptor{},
			path:   "definition.versions[v1].desired",
			target: ErrInvalidVersion,
			reason: ErrorReasonInvalidDesired,
		},
		{
			name:   "desired scalar",
			desc:   stringType(),
			path:   "definition.versions[v1].desired",
			target: ErrInvalidVersion,
			reason: ErrorReasonDesiredNotObject,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateSurfaceLocal(
				tc.desc,
				tc.path,
				ErrorReasonInvalidDesired,
				ErrorReasonDesiredNotObject,
				detailDesiredObjectLikeTemplate,
			)
			requireResourceError(t, err, tc.target, tc.path, tc.reason)
		})
	}
}

func TestValidateSurfaceLocalAcceptsRootRefs(t *testing.T) {
	requireNoError(t, validateSurfaceLocal(
		refType("control.arcoris.dev.WorkerDesired"),
		"definition.versions[v1].desired",
		ErrorReasonInvalidDesired,
		ErrorReasonDesiredNotObject,
		detailDesiredObjectLikeTemplate,
	))
}

func TestValidateSurfaceResolvedAcceptsResolvedObjectRefs(t *testing.T) {
	resolver := fakeResolver{
		types.TypeName("control.arcoris.dev.WorkerDesired"): types.Define(
			"control.arcoris.dev.WorkerDesired",
			types.Object(),
		),
	}

	requireNoError(t, validateSurfaceResolved(
		refType("control.arcoris.dev.WorkerDesired"),
		resolver,
		"definition.versions[v1].desired",
		ErrorReasonInvalidDesired,
		ErrorReasonDesiredNotObject,
		detailDesiredObjectLikeTemplate,
	))
}

func TestValidateSurfaceResolvedRejectsResolvedScalarRefs(t *testing.T) {
	resolver := fakeResolver{
		types.TypeName("control.arcoris.dev.Text"): types.Define(
			"control.arcoris.dev.Text",
			types.String(),
		),
	}

	err := validateSurfaceResolved(
		refType("control.arcoris.dev.Text"),
		resolver,
		"definition.versions[v1].desired",
		ErrorReasonInvalidDesired,
		ErrorReasonDesiredNotObject,
		detailDesiredObjectLikeTemplate,
	)
	requireResourceError(
		t,
		err,
		ErrInvalidVersion,
		"definition.versions[v1].desired",
		ErrorReasonDesiredNotObject,
	)
}
