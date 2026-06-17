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

	metaidentity "arcoris.dev/apimachinery/api/meta/identity"
)

func TestValidateListResult(t *testing.T) {
	request := ListRequest{Resource: validResource(), Scope: MustNamespace("system")}
	tests := []struct {
		name       string
		request    ListRequest
		result     ListResult
		wantReason ErrorReason
		wantCause  error
	}{
		{
			name:    "empty zero revision",
			request: request,
			result:  ListResult{},
		},
		{
			name:    "empty non-zero revision",
			request: request,
			result:  ListResult{Revision: 10},
		},
		{
			name:    "valid item",
			request: request,
			result: ListResult{
				Items:    []ListItem{{Key: validKey(), State: committedStateAt(2, "desired")}},
				Revision: 2,
			},
		},
		{
			name:       "invalid request",
			request:    ListRequest{},
			result:     ListResult{},
			wantReason: ErrorReasonInvalidListResult,
			wantCause:  ErrInvalidListRequest,
		},
		{
			name:    "invalid item",
			request: request,
			result: ListResult{
				Items:    []ListItem{{Key: validKey(), State: State{}}},
				Revision: 2,
			},
			wantReason: ErrorReasonInvalidListItem,
			wantCause:  ErrInvalidRevision,
		},
		{
			name:    "item outside collection",
			request: request,
			result: ListResult{
				Items: []ListItem{{
					Key:   MustKey(validResource(), metaidentity.ObjectName{Namespace: "other", Name: "main"}),
					State: committedStateForObject(1, "other", "main", "desired"),
				}},
				Revision: 1,
			},
			wantReason: ErrorReasonInvalidListItem,
			wantCause:  ErrInvalidListRequest,
		},
		{
			name:    "item newer than result revision",
			request: request,
			result: ListResult{
				Items:    []ListItem{{Key: validKey(), State: committedStateAt(2, "desired")}},
				Revision: 1,
			},
			wantReason: ErrorReasonInvalidListResult,
			wantCause:  ErrInvalidRevision,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateListResult(tt.request, tt.result)
			if tt.wantCause == nil {
				requireNoError(t, err)
				return
			}

			requireErrorIs(t, err, ErrInvalidListResult)
			if !errors.Is(err, tt.wantCause) {
				t.Fatalf("errors.Is(%v, %v) = false", err, tt.wantCause)
			}
			var storeErr *Error
			if !errors.As(err, &storeErr) {
				t.Fatalf("error = %T; want *Error", err)
			}
			if storeErr.Reason != tt.wantReason {
				t.Fatalf("reason = %q; want %q", storeErr.Reason, tt.wantReason)
			}
		})
	}
}
