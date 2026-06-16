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

// TestValidateTermCanonicalizesMetadata verifies leaf validation returns a
// detached canonical term rather than preserving caller-owned value order.
func TestValidateTermCanonicalizesMetadata(t *testing.T) {
	got, err := validateTerm(term{
		kind:           termMetadata,
		metadataDomain: metadataLabels,
		metadataKey:    "env",
		operator:       OperatorIn,
		stringValues:   []string{"qa", "prod", "qa"},
	}, compileOptions{})
	requireNoError(t, err)

	if got.stringValues[0] != "prod" || got.stringValues[1] != "qa" {
		t.Fatalf("stringValues = %#v; want [prod qa]", got.stringValues)
	}
}

// TestValidateTermRejectsUnknownKind keeps malformed private terms out of
// compiled predicates.
func TestValidateTermRejectsUnknownKind(t *testing.T) {
	_, err := validateTerm(term{}, compileOptions{})

	requireErrorIs(t, err, ErrInvalidTerm)
}
