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
)

// TestPredicateQueryRoundTripCanonical verifies Predicate.Query returns a
// detached canonical expression.
func TestPredicateQueryRoundTripCanonical(t *testing.T) {
	left := mustQ(LabelEquals("tier", "backend"))
	right := mustQ(LabelEquals("env", "prod"))
	query := mustAnd(t, left, right, left)

	predicate := mustPredicate(t, query)
	roundTrip := mustPredicate(t, predicate.Query())

	if predicate.Query().expr.canonicalKey() != roundTrip.Query().expr.canonicalKey() {
		t.Fatalf("canonical roundtrip changed")
	}
}

// TestPredicateIsZeroAndPlanAccessors verifies the base Predicate value
// exposes detached canonical query and plan state.
func TestPredicateIsZeroAndPlanAccessors(t *testing.T) {
	if !(Predicate{}).IsZero() {
		t.Fatal("zero Predicate IsZero() = false; want true")
	}

	predicate := mustPredicate(t, mustQ(LabelEquals("env", "prod")))
	if predicate.IsZero() {
		t.Fatal("label predicate IsZero() = true; want false")
	}
	if predicate.Query().IsZero() {
		t.Fatal("label predicate Query().IsZero() = true; want false")
	}
	if len(predicate.Plan().constraints) != 1 {
		t.Fatalf("plan constraints = %d; want 1", len(predicate.Plan().constraints))
	}
}
