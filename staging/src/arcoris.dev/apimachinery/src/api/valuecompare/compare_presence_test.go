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

package valuecompare

import (
	"testing"

	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
)

func TestComparePresenceBothAbsentIsEmpty(t *testing.T) {
	got, done, err := newComparer(Options{}).comparePresence(
		rootField("name"),
		absentOperand(),
		absentOperand(),
		types.String().Type(),
	)
	requireNoError(t, err)

	if !done {
		t.Fatalf("done = false")
	}
	requireResult(t, got, nil, nil, nil)
}

func TestComparePresenceAddedUsesSubtree(t *testing.T) {
	path := rootField("spec")
	descriptor := types.Object(types.Field("image").String().Optional()).Type()
	newValue := value.MustObjectValue(value.ObjectMember("image", value.StringValue("v1")))

	got, done, err := newComparer(Options{}).comparePresence(path, absentOperand(), presentOperand(newValue), descriptor)
	requireNoError(t, err)

	if !done {
		t.Fatalf("done = false")
	}
	requireResult(t, got, paths(path.Field("image")), nil, nil)
}

func TestComparePresenceRemovedUsesSubtree(t *testing.T) {
	path := rootField("spec")
	descriptor := types.Object(types.Field("image").String().Optional()).Type()
	oldValue := value.MustObjectValue(value.ObjectMember("image", value.StringValue("v1")))

	got, done, err := newComparer(Options{}).comparePresence(path, presentOperand(oldValue), absentOperand(), descriptor)
	requireNoError(t, err)

	if !done {
		t.Fatalf("done = false")
	}
	requireResult(t, got, nil, paths(path.Field("image")), nil)
}

func TestComparePresenceBothPresentContinues(t *testing.T) {
	_, done, err := newComparer(Options{}).comparePresence(
		rootField("name"),
		presentOperand(value.StringValue("old")),
		presentOperand(value.StringValue("new")),
		types.String().Type(),
	)
	requireNoError(t, err)

	if done {
		t.Fatalf("done = true")
	}
}
