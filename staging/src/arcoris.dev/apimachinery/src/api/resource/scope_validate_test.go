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

import "testing"

func TestScopeValidate(t *testing.T) {
	requireNoError(t, ScopeGlobal.Validate())
	requireNoError(t, ScopeNamespaced.Validate())

	err := ScopeInvalid.Validate()
	requireResourceError(t, err, ErrInvalidScope, pathScope, ErrorReasonInvalidScope)
	requireDetailContains(t, err, detailScopeSupported)
}
