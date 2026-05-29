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

// Kind identifies an API object kind.
//
// A Kind starts with an uppercase ASCII letter and then contains only ASCII
// letters or digits. It names an API type identity only; object metadata,
// runtime object-kind mutation, and codec-level apiVersion/kind fields are
// intentionally outside this package.
type Kind string

// String returns the canonical kind text without revalidating it.
func (k Kind) String() string {
	return string(k)
}

// IsZero reports whether the kind is absent.
func (k Kind) IsZero() bool {
	return k == ""
}
