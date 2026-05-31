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

	metaidentity "arcoris.dev/apimachinery/api/meta/identity"
	apiobject "arcoris.dev/apimachinery/api/object"
	"arcoris.dev/apimachinery/api/resource"
)

func TestValidateScopeCompatibility(t *testing.T) {
	tests := []struct {
		name      string
		scope     resource.Scope
		namespace string
		wantErr   bool
	}{
		{name: "global empty namespace", scope: resource.ScopeGlobal},
		{name: "global namespaced object", scope: resource.ScopeGlobal, namespace: "system", wantErr: true},
		{name: "namespaced object with namespace", scope: resource.ScopeNamespaced, namespace: "system"},
		{name: "namespaced object without namespace", scope: resource.ScopeNamespaced},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plan := validPlan()
			plan.Resource = resourceDefinition(tt.scope)

			obj := apiobject.NewObserved(
				validTypeMeta("v1"),
				validObjectMeta(metaidentity.Namespace(tt.namespace)),
				testDesired{Replicas: 3},
				testObserved{ReadyReplicas: 2},
			)

			err := Validate(obj, plan)
			if tt.wantErr {
				requireValidationError(
					t,
					err,
					ErrInvalidScope,
					pathObjectNamespace,
					ErrorReasonInvalidScope,
				)
				return
			}

			requireNoError(t, err)
		})
	}
}
