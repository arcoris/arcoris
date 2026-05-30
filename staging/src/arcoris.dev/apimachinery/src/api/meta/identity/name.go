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

// Name is the metadata name of an object.
//
// A name is unique only within resource and scope rules enforced by higher
// layers. This package validates lexical form but never checks uniqueness.
type Name string

// String returns the raw metadata name text.
func (n Name) String() string {
	return string(n)
}

// IsZero reports whether the name is absent.
func (n Name) IsZero() bool {
	return n == ""
}
