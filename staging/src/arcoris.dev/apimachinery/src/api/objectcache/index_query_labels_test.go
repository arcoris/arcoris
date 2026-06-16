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

package objectcache

import (
	"testing"

	"arcoris.dev/apimachinery/api/objectquery"
)

func TestIndexesLabelCandidatesPositiveOperators(t *testing.T) {
	idx := testIndexes()
	tests := []struct {
		name string
		req  objectquery.LabelRequirement
		want []itemRef
	}{
		{
			name: "exists",
			req:  mustLabelExists(t, "tier"),
			want: []itemRef{
				{"system", "worker-1", 1},
				{"system", "worker-2", 2},
				{"other", "worker-3", 3},
			},
		},
		{
			name: "equals",
			req:  mustLabelEquals(t, "env", "qa"),
			want: []itemRef{{"system", "worker-2", 2}},
		},
		{
			name: "in",
			req:  mustLabelIn(t, "env", "prod", "qa"),
			want: []itemRef{
				{"system", "worker-1", 1},
				{"system", "worker-2", 2},
				{"other", "worker-3", 3},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keys, ok := idx.labelCandidates(tt.req)
			if !ok {
				t.Fatal("labelCandidates returned residual-only for positive operator")
			}

			requireKeySetIncludesOnly(t, keys, tt.want...)
		})
	}
}

func TestIndexesLabelCandidatesNegativeOperatorsAreResidual(t *testing.T) {
	idx := testIndexes()
	for _, req := range []objectquery.LabelRequirement{
		mustLabelDoesNotExist(t, "env"),
		mustLabelNotEquals(t, "env", "prod"),
		mustLabelNotIn(t, "env", "prod", "qa"),
	} {
		if _, ok := idx.labelCandidates(req); ok {
			t.Fatalf("labelCandidates(%s) returned indexed candidates for residual operator", req.Operator())
		}
	}
}
