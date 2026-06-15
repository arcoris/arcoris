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

func TestAnnotationRequirementZeroValue(t *testing.T) {
	var req AnnotationRequirement
	if req.req.op != 0 || req.req.key != "" || req.req.values != nil {
		t.Fatalf("zero annotation requirement = %#v; want empty", req)
	}
	if req.Key() != "" {
		t.Fatalf("Key() = %q; want empty", req.Key())
	}
	if req.Operator() != 0 {
		t.Fatalf("Operator() = %s; want zero", req.Operator())
	}
	if req.Values() != nil {
		t.Fatalf("Values() = %#v; want nil", req.Values())
	}
}

func TestAnnotationRequirementAccessors(t *testing.T) {
	req := mustAnnotationIn(t, "note", "qa rollout", "prod rollout", "qa rollout")

	if got := req.Key(); got != "note" {
		t.Fatalf("Key() = %q; want note", got)
	}
	if got := req.Operator(); got != OperatorIn {
		t.Fatalf("Operator() = %s; want %s", got, OperatorIn)
	}
	requireStrings(t, req.Values(), "prod rollout", "qa rollout")
}

func TestAnnotationRequirementValuesDefensiveCopy(t *testing.T) {
	req := mustAnnotationIn(t, "note", "qa rollout", "prod rollout")
	selector := mustAnnotationSelector(t, req)
	values := req.Values()
	values[0] = "mutated"

	requireStrings(t, req.Values(), "prod rollout", "qa rollout")
	if !selector.match(testItem("system", "worker", nil, map[string]string{"note": "prod rollout"})) {
		t.Fatal("mutating Values() copy affected selector match")
	}
}

func mustAnnotationExists(t *testing.T, key string) AnnotationRequirement {
	t.Helper()
	req, err := AnnotationExists(key)
	requireNoError(t, err)

	return req
}

func mustAnnotationDoesNotExist(t *testing.T, key string) AnnotationRequirement {
	t.Helper()
	req, err := AnnotationDoesNotExist(key)
	requireNoError(t, err)

	return req
}

func mustAnnotationEquals(t *testing.T, key string, value string) AnnotationRequirement {
	t.Helper()
	req, err := AnnotationEquals(key, value)
	requireNoError(t, err)

	return req
}

func mustAnnotationNotEquals(t *testing.T, key string, value string) AnnotationRequirement {
	t.Helper()
	req, err := AnnotationNotEquals(key, value)
	requireNoError(t, err)

	return req
}

func mustAnnotationIn(t *testing.T, key string, values ...string) AnnotationRequirement {
	t.Helper()
	req, err := AnnotationIn(key, values...)
	requireNoError(t, err)

	return req
}

func mustAnnotationNotIn(t *testing.T, key string, values ...string) AnnotationRequirement {
	t.Helper()
	req, err := AnnotationNotIn(key, values...)
	requireNoError(t, err)

	return req
}
