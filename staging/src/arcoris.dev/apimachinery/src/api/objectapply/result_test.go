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

package objectapply

import (
	"testing"

	"arcoris.dev/apimachinery/api/fieldownership"
	"arcoris.dev/apimachinery/api/fieldpath"
)

func TestResultZeroValue(t *testing.T) {
	var result Result

	if !result.Object.Desired.IsZero() {
		t.Fatalf("zero Result has non-zero object desired")
	}
	if !result.Ownership.IsEmpty() {
		t.Fatalf("zero Result ownership is not empty")
	}
}

func TestApplyEarlyValidationFailureReturnsZeroResult(t *testing.T) {
	req := testRequest()
	req.Owner = fieldownership.Owner{}

	result, err := Apply(req, Options{})
	requireErrorIs(t, err, ErrInvalidOwner)

	if !result.Object.Desired.IsZero() {
		t.Fatalf("object was populated")
	}
	if !result.Ownership.IsEmpty() {
		t.Fatalf("ownership was populated")
	}
	if !result.Desired.Value.IsZero() {
		t.Fatalf("desired result was populated")
	}
}

func TestApplyDesiredUnsupportedTakeoverReturnsPartialDesiredOnly(t *testing.T) {
	req := testRequest()
	req.Ownership = desiredOwnership(entry("other", fieldpath.Root()))

	result, err := Apply(req, Options{Force: true})
	requireErrorIs(t, err, ErrDesiredApplyFailed)

	requireSet(t, result.Desired.AppliedFields, "$.image")
	if !result.Object.Desired.IsZero() {
		t.Fatalf("object was built")
	}
	if !result.Ownership.IsEmpty() {
		t.Fatalf("ownership was updated")
	}
}
