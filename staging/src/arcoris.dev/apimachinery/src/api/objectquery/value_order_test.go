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
	"time"

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

// TestCompareOrderedDurationUsesSemanticDuration verifies objectquery does not
// fall back to lexicographic duration strings.
func TestCompareOrderedDurationUsesSemanticDuration(t *testing.T) {
	cmp, ok := compareOrdered(value.DurationValue(10*time.Second), value.DurationValue(2*time.Second))
	if !ok || cmp <= 0 {
		t.Fatalf("duration compare = (%d, %v); want greater/true", cmp, ok)
	}
}
