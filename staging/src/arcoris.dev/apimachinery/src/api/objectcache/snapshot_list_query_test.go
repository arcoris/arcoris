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

func TestSnapshotListQueriesMatchObjectQueryFullScan(t *testing.T) {
	source := testItems()
	snapshot := mustSnapshot(t, testListResult(14, source...))
	tests := []struct {
		name  string
		query objectquery.Query
	}{
		{
			name: "identity namespace",
			query: objectquery.Query{
				Identity: objectquery.IdentitySelector{
					Namespace: mustNamespaceEquals(t, "system"),
				},
			},
		},
		{
			name: "identity name",
			query: objectquery.Query{
				Identity: objectquery.IdentitySelector{
					Name: mustNameEquals(t, "worker-3"),
				},
			},
		},
		{
			name: "label exists",
			query: objectquery.Query{
				Labels: mustLabelSelector(t, mustLabelExists(t, "env")),
			},
		},
		{
			name: "label equals",
			query: objectquery.Query{
				Labels: mustLabelSelector(t, mustLabelEquals(t, "env", "prod")),
			},
		},
		{
			name: "label in",
			query: objectquery.Query{
				Labels: mustLabelSelector(t, mustLabelIn(t, "env", "prod", "qa")),
			},
		},
		{
			name: "label does not exist",
			query: objectquery.Query{
				Labels: mustLabelSelector(t, mustLabelDoesNotExist(t, "env")),
			},
		},
		{
			name: "label not equals",
			query: objectquery.Query{
				Labels: mustLabelSelector(t, mustLabelNotEquals(t, "env", "prod")),
			},
		},
		{
			name: "label not in",
			query: objectquery.Query{
				Labels: mustLabelSelector(t, mustLabelNotIn(t, "env", "prod", "qa")),
			},
		},
		{
			name: "annotation exists",
			query: objectquery.Query{
				Annotations: mustAnnotationSelector(t, mustAnnotationExists(t, "team")),
			},
		},
		{
			name: "annotation equals",
			query: objectquery.Query{
				Annotations: mustAnnotationSelector(t, mustAnnotationEquals(t, "team", "core")),
			},
		},
		{
			name: "annotation in",
			query: objectquery.Query{
				Annotations: mustAnnotationSelector(t, mustAnnotationIn(t, "team", "core", "tools")),
			},
		},
		{
			name: "annotation does not exist",
			query: objectquery.Query{
				Annotations: mustAnnotationSelector(t, mustAnnotationDoesNotExist(t, "team")),
			},
		},
		{
			name: "annotation not equals",
			query: objectquery.Query{
				Annotations: mustAnnotationSelector(t, mustAnnotationNotEquals(t, "team", "core")),
			},
		},
		{
			name: "annotation not in",
			query: objectquery.Query{
				Annotations: mustAnnotationSelector(t, mustAnnotationNotIn(t, "team", "core", "tools")),
			},
		},
		{
			name: "combined",
			query: objectquery.Query{
				Identity: objectquery.IdentitySelector{
					Namespace: mustNamespaceEquals(t, "system"),
				},
				Labels: mustLabelSelector(
					t,
					mustLabelEquals(t, "tier", "backend"),
					mustLabelIn(t, "env", "prod", "qa"),
				),
				Annotations: mustAnnotationSelector(
					t,
					mustAnnotationEquals(t, "zone", "east"),
				),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertSnapshotListMatchesObjectQueryFullScan(t, snapshot, source, tt.query)
		})
	}
}
