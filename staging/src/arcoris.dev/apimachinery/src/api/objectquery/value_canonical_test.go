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
	"strings"
	"testing"

	"arcoris.dev/apimachinery/api/value"
)

// TestCanonicalValueKeyDistinguishesKinds verifies literal equality keys do
// not collapse different concrete value kinds with similar text.
func TestCanonicalValueKeyDistinguishesKinds(t *testing.T) {
	stringKey := canonicalValueKey(value.StringValue("1"))
	intKey := canonicalValueKey(value.Int64Value(1))

	if stringKey == intKey {
		t.Fatalf("string and integer canonical keys are equal: %q", stringKey)
	}
	if !strings.HasPrefix(canonicalValueKey(value.NullValue()), "null\x00") {
		t.Fatal("null canonical key lost null prefix")
	}
}

// TestCanonicalValueKeyHandlesCompositeValues verifies records and lists have
// deterministic recursive keys.
func TestCanonicalValueKeyHandlesCompositeValues(t *testing.T) {
	record := value.MustRecordValue(value.MustRecordMember("name", value.StringValue("api")))
	list := value.MustListValue(value.StringValue("api"), value.Int64Value(1))

	if !strings.HasPrefix(canonicalValueKey(record), "record\x00") {
		t.Fatalf("record canonical key = %q", canonicalValueKey(record))
	}
	if !strings.HasPrefix(canonicalValueKey(list), "list\x00") {
		t.Fatalf("list canonical key = %q", canonicalValueKey(list))
	}
}

// TestCanonicalValueKeyAvoidsDelimiterCollisions verifies structural keys do
// not collide when user strings contain punctuation used by older encodings.
func TestCanonicalValueKeyAvoidsDelimiterCollisions(t *testing.T) {
	left := value.MustRecordValue(value.MustRecordMember("a,b", value.StringValue("c=d")))
	right := value.MustRecordValue(value.MustRecordMember("a", value.StringValue("b,c=d")))

	if canonicalValueKey(left) == canonicalValueKey(right) {
		t.Fatalf("canonical keys collided: %q", canonicalValueKey(left))
	}
}
