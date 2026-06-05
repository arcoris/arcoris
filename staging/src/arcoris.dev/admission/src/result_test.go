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

package admission

import (
	"reflect"
	"testing"
)

func TestResultPresenceFlags(t *testing.T) {
	t.Parallel()

	result := resultWith(GrantDecision(ReasonAdmitted), "grant", true, "metadata", true)
	if !result.HasGrant() {
		t.Fatal("HasGrant() = false, want true")
	}
	if !result.HasMetadata() {
		t.Fatal("HasMetadata() = false, want true")
	}
}

func TestResultDoesNotExposeDecisionProjectionMethods(t *testing.T) {
	t.Parallel()

	resultType := reflect.TypeOf(Result[string, string]{})
	for _, name := range []string{"IsAdmitted", "IsDenied", "IsQueued", "IsDeferred", "HasSideEffect"} {
		requireNoMethod(t, resultType, name)
	}
}
