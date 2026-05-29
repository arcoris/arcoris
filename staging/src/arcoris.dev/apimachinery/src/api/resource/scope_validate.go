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

// Validate checks whether s identifies a supported resource scope.
//
// Validation is lexical and descriptor-local. It does not define namespace
// objects, authorization policy, routing behavior, or storage partitioning.
func (s Scope) Validate() error {
	if s.IsValid() {
		return nil
	}
	return scopeError(ErrorReasonInvalidScope, detailScopeSupported)
}
