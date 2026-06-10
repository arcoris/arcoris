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

package fieldownership

import "testing"

func TestOwnerString(t *testing.T) {
	requireEqual(t, owner("user-cli").String(), "user-cli")
}

func TestOwnerIsZero(t *testing.T) {
	requireEqual(t, Owner{}.IsZero(), true)
	requireEqual(t, owner("user-cli").IsZero(), false)
}

func TestOwnerCompare(t *testing.T) {
	requireEqual(t, owner("a").Compare(owner("b")), -1)
	requireEqual(t, owner("b").Compare(owner("a")), 1)
	requireEqual(t, owner("a").Compare(owner("a")), 0)
}

func TestNewOwnerAcceptsValidExamples(t *testing.T) {
	for _, value := range []string{
		"user-cli",
		"terraform",
		"status-controller",
		"arcoris.dev/controller",
		"user:anton",
	} {
		owner, err := NewOwner(value)

		requireNoError(t, err)
		requireEqual(t, owner.String(), value)
	}
}

func TestMustOwnerPanicsOnInvalidOwner(t *testing.T) {
	requirePanic(t, func() {
		MustOwner("")
	})
}
