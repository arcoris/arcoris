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

import "fmt"

// String returns stable diagnostic text for k.
//
// The string form is useful for diagnostics and deterministic in-memory
// hashing. It is not a transport path, authorization resource, or serialized
// API identity.
func (k Key) String() string {
	if k.Resource.IsZero() && k.Object.IsZero() {
		return ""
	}

	return fmt.Sprintf("%s/%s", k.Resource, k.Object)
}
