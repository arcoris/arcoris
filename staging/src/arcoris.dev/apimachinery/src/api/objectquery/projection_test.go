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

	"arcoris.dev/apimachinery/api/objectstore"
)

// TestProjectChange verifies every membership transition kind.
func TestProjectChange(t *testing.T) {
	predicate := mustPredicate(t, mustQ(LabelEquals("env", "prod")))
	before := testItem("system", "worker-1", 1, map[string]string{"env": "qa"}, nil, desiredRecord("api", "qa", 1))
	after := testItem("system", "worker-1", 2, map[string]string{"env": "prod"}, nil, desiredRecord("api", "prod", 1))
	other := testItem("system", "worker-2", 3, map[string]string{"env": "dev"}, nil, desiredRecord("api", "dev", 1))

	tests := []struct {
		name string
		ch   objectstore.Change
		want ChangeProjectionKind
	}{
		{name: "create matching", ch: objectstore.MustCreatedChange(after.Key, after.State), want: ChangeProjectionEntered},
		{name: "create ignored", ch: objectstore.MustCreatedChange(other.Key, other.State), want: ChangeProjectionIgnored},
		{name: "update entered", ch: objectstore.MustUpdatedChange(before.Key, before.State, after.State), want: ChangeProjectionEntered},
		{name: "update updated", ch: objectstore.MustUpdatedChange(after.Key, after.State, testItem("system", "worker-1", 4, map[string]string{"env": "prod"}, nil, desiredRecord("api", "prod", 1)).State), want: ChangeProjectionUpdated},
		{name: "update left", ch: objectstore.MustUpdatedChange(after.Key, after.State, other.State), want: ChangeProjectionLeft},
		{name: "update ignored", ch: objectstore.MustUpdatedChange(before.Key, before.State, other.State), want: ChangeProjectionIgnored},
		{name: "delete left", ch: objectstore.MustDeletedChange(after.Key, after.State, 5), want: ChangeProjectionLeft},
		{name: "delete ignored", ch: objectstore.MustDeletedChange(other.Key, other.State, 6), want: ChangeProjectionIgnored},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := predicate.ProjectChange(tt.ch)
			requireNoError(t, err)
			if got.Kind != tt.want {
				t.Fatalf("projection = %v; want %v", got.Kind, tt.want)
			}
		})
	}
}

// TestProjectChangeInvalidChange verifies objectstore.Change validation is
// preserved by projection.
func TestProjectChangeInvalidChange(t *testing.T) {
	_, err := mustPredicate(t, All()).ProjectChange(objectstore.Change{})

	requireErrorIs(t, err, ErrInvalidChange)
	requireStructuredQueryError(t, err, ErrorReasonInvalidChange, "query.change")
}
