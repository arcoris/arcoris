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

package stamp

// ResourceVersion is an opaque API/storage concurrency and change token.
//
// It is not globally comparable, not a number, and not snapshot.Revision.
type ResourceVersion string

// String returns the opaque resource version text.
func (v ResourceVersion) String() string {
	return string(v)
}

// IsZero reports whether the token is absent.
func (v ResourceVersion) IsZero() bool {
	return v == ""
}
