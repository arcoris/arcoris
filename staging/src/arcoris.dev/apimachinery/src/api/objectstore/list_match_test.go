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
	"testing"

	apiidentity "arcoris.dev/apimachinery/api/identity"
	metaidentity "arcoris.dev/apimachinery/api/meta/identity"
)

func TestKeyMatchesListRequest(t *testing.T) {
	otherResource := apiidentity.GroupVersionResource{
		Group:    "control.arcoris.dev",
		Version:  "v1",
		Resource: "tasks",
	}
	key := validKey()

	tests := []struct {
		name    string
		key     Key
		request ListRequest
		want    bool
	}{
		{
			name:    "matching resource all namespaces",
			key:     key,
			request: ListRequest{Resource: validResource(), Scope: AllNamespaces()},
			want:    true,
		},
		{
			name:    "non-matching resource",
			key:     MustKey(otherResource, validObjectName()),
			request: ListRequest{Resource: validResource(), Scope: AllNamespaces()},
		},
		{
			name:    "matching namespace",
			key:     key,
			request: ListRequest{Resource: validResource(), Scope: MustNamespace("system")},
			want:    true,
		},
		{
			name: "non-matching namespace",
			key: MustKey(validResource(), metaidentity.ObjectName{
				Namespace: "other",
				Name:      "main",
			}),
			request: ListRequest{Resource: validResource(), Scope: MustNamespace("system")},
		},
		{
			name:    "zero request",
			key:     key,
			request: ListRequest{},
		},
		{
			name:    "zero key",
			key:     Key{},
			request: ListRequest{Resource: validResource(), Scope: AllNamespaces()},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := KeyMatchesListRequest(tt.key, tt.request); got != tt.want {
				t.Fatalf("KeyMatchesListRequest() = %v; want %v", got, tt.want)
			}
		})
	}
}

func TestChangeMatchesListRequest(t *testing.T) {
	request := ListRequest{Resource: validResource(), Scope: MustNamespace("system")}
	tests := []struct {
		name   string
		change Change
		want   bool
	}{
		{name: "create", change: MustCreatedChange(validKey(), committedStateAt(1, "created")), want: true},
		{name: "update", change: MustUpdatedChange(validKey(), committedStateAt(1, "before"), committedStateAt(2, "after")), want: true},
		{name: "delete", change: MustDeletedChange(validKey(), committedStateAt(1, "before"), 2), want: true},
		{name: "invalid change", change: Change{}},
		{
			name: "outside namespace",
			change: MustCreatedChange(
				MustKey(validResource(), metaidentity.ObjectName{Namespace: "other", Name: "main"}),
				committedStateForObject(1, "other", "main", "created"),
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ChangeMatchesListRequest(tt.change, request); got != tt.want {
				t.Fatalf("ChangeMatchesListRequest() = %v; want %v", got, tt.want)
			}
		})
	}
}
