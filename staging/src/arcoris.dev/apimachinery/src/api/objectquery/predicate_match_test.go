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

import "testing"

func TestPredicateMatchCombinesSectionsWithAnd(t *testing.T) {
	predicate, err := Compile(Query{
		Identity: mustWithObject(t, "system", "worker"),
		Labels: mustLabelSelector(
			t,
			mustLabelEquals(t, "env", "prod"),
		),
		Annotations: mustAnnotationSelector(
			t,
			mustAnnotationEquals(t, "note", "rollout"),
		),
	})
	requireNoError(t, err)

	tests := []struct {
		name  string
		item  objectNameAndMetadata
		match bool
	}{
		{
			name: "all match",
			item: objectNameAndMetadata{
				namespace:   "system",
				name:        "worker",
				labels:      map[string]string{"env": "prod"},
				annotations: map[string]string{"note": "rollout"},
			},
			match: true,
		},
		{
			name: "namespace mismatch",
			item: objectNameAndMetadata{
				namespace:   "other",
				name:        "worker",
				labels:      map[string]string{"env": "prod"},
				annotations: map[string]string{"note": "rollout"},
			},
		},
		{
			name: "name mismatch",
			item: objectNameAndMetadata{
				namespace:   "system",
				name:        "other",
				labels:      map[string]string{"env": "prod"},
				annotations: map[string]string{"note": "rollout"},
			},
		},
		{
			name: "label mismatch",
			item: objectNameAndMetadata{
				namespace:   "system",
				name:        "worker",
				labels:      map[string]string{"env": "qa"},
				annotations: map[string]string{"note": "rollout"},
			},
		},
		{
			name: "annotation mismatch",
			item: objectNameAndMetadata{
				namespace:   "system",
				name:        "worker",
				labels:      map[string]string{"env": "prod"},
				annotations: map[string]string{"note": "other"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item := testItem(tt.item.namespace, tt.item.name, tt.item.labels, tt.item.annotations)
			if got := predicate.Match(item); got != tt.match {
				t.Fatalf("Match() = %v; want %v", got, tt.match)
			}
		})
	}
}

type objectNameAndMetadata struct {
	namespace   string
	name        string
	labels      map[string]string
	annotations map[string]string
}
