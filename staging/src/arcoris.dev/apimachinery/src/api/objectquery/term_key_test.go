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

package objectquery

import (
	"testing"

	apiidentity "arcoris.dev/apimachinery/api/identity"
	metaidentity "arcoris.dev/apimachinery/api/meta/identity"
)

// TestKeyAndResourceTerms verifies storage-key based term matching.
func TestKeyAndResourceTerms(t *testing.T) {
	items := testItems()
	tests := []struct {
		name  string
		query Query
		want  []metaidentity.Name
	}{
		{
			name:  "resource",
			query: mustQ(ResourceEquals(testResource)),
			want:  []metaidentity.Name{"worker-1", "worker-2", "worker-3", "worker-4"},
		},
		{
			name:  "namespace",
			query: mustQ(ObjectInNamespace("other")),
			want:  []metaidentity.Name{"worker-3"},
		},
		{
			name:  "name",
			query: mustQ(ObjectWithName("worker-2")),
			want:  []metaidentity.Name{"worker-2"},
		},
		{
			name:  "object",
			query: mustQ(ObjectEquals("system", "worker-1")),
			want:  []metaidentity.Name{"worker-1"},
		},
		{
			name:  "key",
			query: mustQ(KeyEquals(items[2].Key)),
			want:  []metaidentity.Name{"worker-3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requireNames(t, mustPredicate(t, tt.query).Filter(items), tt.want...)
		})
	}
}

// TestKeyTermValidationErrors verifies lexical and resource validation errors.
func TestKeyTermValidationErrors(t *testing.T) {
	_, err := ResourceEquals(apiidentity.GroupVersionResource{})
	requireErrorIs(t, err, ErrInvalidTerm)

	_, err = ObjectWithName("")
	requireErrorIs(t, err, ErrInvalidTerm)

	_, err = ObjectInNamespace("Bad")
	requireErrorIs(t, err, ErrInvalidTerm)
}
