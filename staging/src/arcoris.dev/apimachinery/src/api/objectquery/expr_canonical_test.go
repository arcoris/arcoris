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

// TestCanonicalExprKeyIsStructural verifies expression keys encode boolean
// shape and are not just leaf-term text.
func TestCanonicalExprKeyIsStructural(t *testing.T) {
	label := mustQ(LabelEquals("env", "prod")).expr
	notLabel := mustNot(t, queryFromExpr(label)).expr

	if label.canonicalKey() == notLabel.canonicalKey() {
		t.Fatalf("term and not-term canonical keys are equal: %q", label.canonicalKey())
	}
	if got := canonicalExprKey(nil); got != "all" {
		t.Fatalf("canonicalExprKey(nil) = %q; want all", got)
	}
}

// TestJoinExprKeysPreservesCallerOrder verifies the helper does not sort by
// itself; sorting belongs to normalization.
func TestJoinExprKeysPreservesCallerOrder(t *testing.T) {
	left := mustQ(LabelEquals("tier", "backend")).expr
	right := mustQ(LabelEquals("env", "prod")).expr

	got := joinExprKeys([]*expr{left, right})
	want := left.canonicalKey() + "," + right.canonicalKey()
	if got != want {
		t.Fatalf("joinExprKeys = %q; want %q", got, want)
	}
}
