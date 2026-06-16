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

// TestTermCloneDetachesLiteralSlices verifies leaf terms are immutable after
// they enter Query or Predicate values.
func TestTermCloneDetachesLiteralSlices(t *testing.T) {
	source := term{
		kind:           termMetadata,
		metadataDomain: metadataLabels,
		metadataKey:    "env",
		operator:       OperatorIn,
		stringValues:   []string{"prod", "qa"},
		values:         []value.Value{value.StringValue("prod")},
	}

	clone := source.clone()
	source.stringValues[0] = "dev"
	source.values[0] = value.StringValue("dev")

	if clone.stringValues[0] != "prod" {
		t.Fatalf("cloned string value = %q; want prod", clone.stringValues[0])
	}
	got, _ := clone.values[0].AsString()
	if got != "prod" {
		t.Fatalf("cloned value = %q; want prod", got)
	}
}
