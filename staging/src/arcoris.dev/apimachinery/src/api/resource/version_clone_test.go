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

package resource

import (
	"testing"

	"arcoris.dev/apimachinery/api/identity"
)

func TestCloneVersions(t *testing.T) {
	input := []VersionDefinition{validVersion()}
	cloned := cloneVersions(input)
	cloned[0] = NewVersion(identity.Version("v2"), objectType())

	requireEqual(t, input[0].Version(), identity.Version("v1"))
}

func TestCloneVersionsEmpty(t *testing.T) {
	if cloneVersions(nil) != nil {
		t.Fatalf("cloneVersions(nil) must return nil")
	}

	if cloneVersions([]VersionDefinition{}) != nil {
		t.Fatalf("cloneVersions(empty) must return nil")
	}
}
