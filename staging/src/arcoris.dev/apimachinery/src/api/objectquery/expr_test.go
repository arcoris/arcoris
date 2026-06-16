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

// TestExprCloneDetachesChildrenAndTerms verifies private expression copies do
// not share mutable literal slices with their source.
func TestExprCloneDetachesChildrenAndTerms(t *testing.T) {
	source := mustQ(LabelIn("env", "prod", "qa")).expr
	clone := source.clone()

	source.term.stringValues[0] = "dev"

	if clone.term.stringValues[0] != "prod" {
		t.Fatalf("clone stringValues[0] = %q; want prod", clone.term.stringValues[0])
	}
}

// TestExprCanonicalKeyTreatsNilAsAll locks down the nil expression sentinel
// used by zero Query and zero Predicate values.
func TestExprCanonicalKeyTreatsNilAsAll(t *testing.T) {
	var e *expr

	if got := e.canonicalKey(); got != "all" {
		t.Fatalf("nil canonical key = %q; want all", got)
	}
}
