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

func TestValidateListRequest(t *testing.T) {
	tests := []struct {
		name   string
		req    ListRequest
		target error
		reason ErrorReason
	}{
		{
			name: "all namespaces",
			req:  ListRequest{Resource: validResource(), Scope: AllNamespaces()},
		},
		{
			name: "namespace",
			req:  ListRequest{Resource: validResource(), Scope: MustNamespace("system")},
		},
		{
			name:   "zero scope",
			req:    ListRequest{Resource: validResource()},
			target: ErrInvalidListRequest,
			reason: ErrorReasonInvalidListScope,
		},
		{
			name:   "unknown scope kind",
			req:    ListRequest{Resource: validResource(), Scope: ListScope{kind: ListScopeKind(99)}},
			target: ErrInvalidListRequest,
			reason: ErrorReasonInvalidListScope,
		},
		{
			name:   "namespace scope without namespace",
			req:    ListRequest{Resource: validResource(), Scope: ListScope{kind: ListScopeNamespace}},
			target: ErrInvalidListRequest,
			reason: ErrorReasonInvalidListScope,
		},
		{
			name:   "namespace scope with invalid namespace",
			req:    ListRequest{Resource: validResource(), Scope: ListScope{kind: ListScopeNamespace, namespace: "System"}},
			target: ErrInvalidListRequest,
			reason: ErrorReasonInvalidListScope,
		},
		{
			name:   "all namespaces carries namespace",
			req:    ListRequest{Resource: validResource(), Scope: ListScope{kind: ListScopeAll, namespace: "system"}},
			target: ErrInvalidListRequest,
			reason: ErrorReasonInvalidListScope,
		},
		{
			name:   "invalid resource",
			req:    ListRequest{Scope: AllNamespaces()},
			target: ErrInvalidListRequest,
			reason: ErrorReasonInvalidListRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateListRequest(tt.req)
			if tt.target == nil {
				requireNoError(t, err)
				return
			}

			requireErrorIs(t, err, tt.target)
			var storeErr *Error
			if !errors.As(err, &storeErr) {
				t.Fatalf("error = %T; want *Error", err)
			}
			if storeErr.Reason != tt.reason {
				t.Fatalf("reason = %q; want %q", storeErr.Reason, tt.reason)
			}
		})
	}
}

func TestInNamespaceRejectsInvalidNamespace(t *testing.T) {
	_, err := InNamespace("System")

	requireErrorIs(t, err, ErrInvalidListRequest)
}
