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

package identity

import "testing"

func TestGroupZeroIsCoreGroup(t *testing.T) {
	if !CoreGroup.IsZero() {
		t.Fatalf("CoreGroup should be zero")
	}

	requireNoError(t, CoreGroup.Validate())
}

func TestGroupStructuredDiagnostics(t *testing.T) {
	err := Group("apps").Validate()
	requireIdentityError(t, err, identityNameGroup, ErrorReasonInvalidForm)
	requireDetailContains(t, err, "qualified DNS subdomain")
}
