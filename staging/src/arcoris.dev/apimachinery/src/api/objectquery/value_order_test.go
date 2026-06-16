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

	"arcoris.dev/apimachinery/api/value"
)

// TestCompareOrdered verifies ordered value comparison accepts supported kinds
// and rejects unsupported or mismatched kinds.
func TestCompareOrdered(t *testing.T) {
	cmp, ok := compareOrdered(value.Int64Value(1), value.Int64Value(2))
	if !ok || cmp >= 0 {
		t.Fatalf("integer compare = (%d, %v); want less/true", cmp, ok)
	}
	if _, ok := compareOrdered(value.StringValue("a"), value.StringValue("b")); ok {
		t.Fatal("string compareOrdered ok = true; want false")
	}
	if _, ok := compareOrdered(value.Int64Value(1), value.StringValue("1")); ok {
		t.Fatal("mismatched compareOrdered ok = true; want false")
	}
}

// TestCompareCanonicalKeys verifies the fallback comparator has stable
// three-way semantics.
func TestCompareCanonicalKeys(t *testing.T) {
	if compareCanonicalKeys(value.StringValue("a"), value.StringValue("b")) >= 0 {
		t.Fatal("canonical a < b failed")
	}
	if compareCanonicalKeys(value.StringValue("b"), value.StringValue("a")) <= 0 {
		t.Fatal("canonical b > a failed")
	}
	if compareCanonicalKeys(value.StringValue("a"), value.StringValue("a")) != 0 {
		t.Fatal("canonical a == a failed")
	}
}
