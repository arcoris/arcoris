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

// TestCanonicalValuesSortsDeduplicatesAndRejectsZero verifies field literal
// sets are canonical and intentional.
func TestCanonicalValuesSortsDeduplicatesAndRejectsZero(t *testing.T) {
	values, err := canonicalValues([]value.Value{
		value.StringValue("qa"),
		value.StringValue("prod"),
		value.StringValue("qa"),
	})
	requireNoError(t, err)

	if len(values) != 2 {
		t.Fatalf("canonical values = %d; want 2", len(values))
	}
	first, _ := values[0].AsString()
	second, _ := values[1].AsString()
	if first != "prod" || second != "qa" {
		t.Fatalf("canonical order = [%s %s]; want [prod qa]", first, second)
	}

	_, err = canonicalValues([]value.Value{{}})
	requireErrorIs(t, err, ErrInvalidTerm)
}
