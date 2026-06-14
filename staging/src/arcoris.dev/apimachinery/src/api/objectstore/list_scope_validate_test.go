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

package objectstore

import (
	"errors"
	"testing"
)

func TestValidateListScope(t *testing.T) {
	tests := []struct {
		name  string
		scope ListScope
		valid bool
	}{
		{name: "all namespaces", scope: AllNamespaces(), valid: true},
		{name: "namespace", scope: MustNamespace("system"), valid: true},
		{name: "zero", scope: ListScope{}},
		{name: "unknown kind", scope: ListScope{kind: ListScopeKind(99)}},
		{name: "namespace without value", scope: ListScope{kind: ListScopeNamespace}},
		{name: "namespace with invalid value", scope: ListScope{kind: ListScopeNamespace, namespace: "System"}},
		{name: "all namespaces with value", scope: ListScope{kind: ListScopeAll, namespace: "system"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateListScope(tt.scope)
			if tt.valid {
				requireNoError(t, err)
				if !tt.scope.IsValid() {
					t.Fatalf("IsValid() = false; want true")
				}
				return
			}

			requireErrorIs(t, err, ErrInvalidListRequest)
			if tt.scope.IsValid() {
				t.Fatalf("IsValid() = true; want false")
			}

			var storeErr *Error
			if !errors.As(err, &storeErr) {
				t.Fatalf("error = %T; want *Error", err)
			}
			if storeErr.Reason != ErrorReasonInvalidListScope {
				t.Fatalf("reason = %q; want %q", storeErr.Reason, ErrorReasonInvalidListScope)
			}
		})
	}
}

func TestInNamespaceRejectsInvalidNamespace(t *testing.T) {
	_, err := InNamespace("System")

	requireErrorIs(t, err, ErrInvalidListRequest)
}
