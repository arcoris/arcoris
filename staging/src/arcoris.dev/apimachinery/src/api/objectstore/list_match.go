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

package objectstore

// KeyMatchesListRequest reports whether key belongs to request's structural collection.
//
// The check is intentionally storage-level only: resource equality plus
// ListScope namespace matching. It does not apply objectquery predicates,
// descriptor scope rules, authorization, pagination, or lifecycle policy.
func KeyMatchesListRequest(key Key, request ListRequest) bool {
	if err := ValidateKey(key); err != nil {
		return false
	}
	if err := ValidateListRequest(request); err != nil {
		return false
	}
	if key.Resource != request.Resource {
		return false
	}

	return keyMatchesListScope(key, request.Scope)
}

// ChangeMatchesListRequest reports whether change belongs to request's structural collection.
//
// The change must be a structurally valid committed transition and its Key must
// match request. This helper does not interpret selectors, watch streams,
// caches, descriptors, or transport-specific request policy.
func ChangeMatchesListRequest(change Change, request ListRequest) bool {
	if err := ValidateChange(change); err != nil {
		return false
	}

	return KeyMatchesListRequest(change.Key, request)
}

// keyMatchesListScope applies only the namespace portion of a validated ListScope.
func keyMatchesListScope(key Key, scope ListScope) bool {
	if scope.IsAllNamespaces() {
		return true
	}
	if scope.IsNamespace() {
		return key.Object.Namespace == scope.Namespace()
	}

	return false
}
