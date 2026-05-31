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

package objectvalidation

import (
	"testing"

	apiidentity "arcoris.dev/apimachinery/api/identity"
	"arcoris.dev/apimachinery/api/resource"
)

func TestValidateChecksResourceGroupAndKind(t *testing.T) {
	tests := []struct {
		name     string
		resource resource.Definition
		wantErr  bool
	}{
		{
			name:     "matching",
			resource: resourceDefinition(resource.ScopeNamespaced),
		},
		{
			name: "wrong group",
			resource: resource.NewDefinition(
				apiidentity.Group("other.arcoris.dev"),
				apiidentity.Kind("Worker"),
				apiidentity.Resource("workers"),
				resource.ScopeNamespaced,
				versionWithObserved("v1"),
			),
			wantErr: true,
		},
		{
			name: "wrong kind",
			resource: resource.NewDefinition(
				apiidentity.Group("control.arcoris.dev"),
				apiidentity.Kind("Other"),
				apiidentity.Resource("workers"),
				resource.ScopeNamespaced,
				versionWithObserved("v1"),
			),
			wantErr: true,
		},
		{
			name: "resource collection name ignored",
			resource: resource.NewDefinition(
				apiidentity.Group("control.arcoris.dev"),
				apiidentity.Kind("Worker"),
				apiidentity.Resource("workerobjects"),
				resource.ScopeNamespaced,
				versionWithObserved("v1"),
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plan := validPlan()
			plan.Resource = tt.resource

			err := Validate(validObject(), plan)
			if tt.wantErr {
				requireValidationError(
					t,
					err,
					ErrResourceMismatch,
					pathObjectTypeMeta,
					ErrorReasonResourceMismatch,
				)
				return
			}

			requireNoError(t, err)
		})
	}
}
