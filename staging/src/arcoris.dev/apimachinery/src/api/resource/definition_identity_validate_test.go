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

	"arcoris.dev/apimachinery/api/identity"
)

func TestValidateDefinitionIdentityRejectsInvalidIdentity(t *testing.T) {
	cases := []struct {
		name   string
		def    Definition
		path   string
		reason ErrorReason
	}{
		{
			name: "invalid group",
			def: NewDefinition(
				identity.Group("apps"),
				identity.Kind("Worker"),
				identity.Resource("workers"),
				ScopeNamespaced,
				validVersion(),
			),
			path:   pathDefinitionGroup,
			reason: ErrorReasonInvalidGroup,
		},
		{
			name: "invalid kind",
			def: NewDefinition(
				identity.Group("control.arcoris.dev"),
				identity.Kind("worker"),
				identity.Resource("workers"),
				ScopeNamespaced,
				validVersion(),
			),
			path:   pathDefinitionKind,
			reason: ErrorReasonInvalidKind,
		},
		{
			name: "invalid resource",
			def: NewDefinition(
				identity.Group("control.arcoris.dev"),
				identity.Kind("Worker"),
				identity.Resource("Workers"),
				ScopeNamespaced,
				validVersion(),
			),
			path:   pathDefinitionResource,
			reason: ErrorReasonInvalidResource,
		},
		{
			name: "invalid scope",
			def: NewDefinition(
				identity.Group("control.arcoris.dev"),
				identity.Kind("Worker"),
				identity.Resource("workers"),
				ScopeInvalid,
				validVersion(),
			),
			path:   pathDefinitionScope,
			reason: ErrorReasonInvalidScope,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateDefinitionIdentity(tc.def)
			requireResourceError(t, err, ErrInvalidDefinition, tc.path, tc.reason)
		})
	}
}
