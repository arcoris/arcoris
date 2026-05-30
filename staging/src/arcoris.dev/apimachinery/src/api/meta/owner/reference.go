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

package owner

import metaidentity "arcoris.dev/apimachinery/api/meta/identity"

// Reference is one owner reference metadata entry.
type Reference struct {
	// Ref is the object reference stored as ownership metadata.
	Ref metaidentity.ObjectReference `json:"ref"`
	// Controller marks the single controlling owner when true.
	Controller bool `json:"controller,omitempty"`
}

// IsZero reports whether the owner reference entry is absent.
func (r Reference) IsZero() bool {
	return r.Ref.IsZero() && !r.Controller
}
