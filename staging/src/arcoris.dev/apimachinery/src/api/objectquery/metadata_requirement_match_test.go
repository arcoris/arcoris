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

func TestMetadataRequirementMatchUnknownOperatorDoesNotMatch(t *testing.T) {
	req := metadataRequirement{key: "env", op: Operator(255)}
	if req.match(func(string) (string, bool) { return "prod", true }) {
		t.Fatal("unknown operator matched")
	}
}

func TestMetadataRequirementMatchNegativeOperatorsMatchAbsentKeys(t *testing.T) {
	lookup := func(string) (string, bool) { return "", false }

	notEquals := metadataRequirement{key: "env", op: OperatorNotEquals, values: []string{"prod"}}
	if !notEquals.match(lookup) {
		t.Fatal("notEquals did not match absent key")
	}

	notIn := metadataRequirement{key: "env", op: OperatorNotIn, values: []string{"prod"}}
	if !notIn.match(lookup) {
		t.Fatal("notIn did not match absent key")
	}
}

func TestMetadataRequirementMatchMembershipUsesCanonicalSortedValues(t *testing.T) {
	lookup := func(string) (string, bool) { return "qa", true }

	in := metadataRequirement{key: "env", op: OperatorIn, values: []string{"prod", "qa"}}
	if !in.match(lookup) {
		t.Fatal("in did not match canonical sorted value set")
	}

	notIn := metadataRequirement{key: "env", op: OperatorNotIn, values: []string{"prod", "qa"}}
	if notIn.match(lookup) {
		t.Fatal("notIn matched value present in canonical sorted value set")
	}
}
