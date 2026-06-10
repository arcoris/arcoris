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

// UID identifies a concrete object incarnation.
//
// UID is opaque metadata assigned by higher layers. It is not a name, resource
// version, generation, timestamp, or UUID parser.
type UID string

// String returns the raw UID text without validating it.
func (u UID) String() string {
	return string(u)
}

// CanonicalText validates the UID and returns its canonical text.
func (u UID) CanonicalText() (string, error) {
	if err := u.ValidateLexical(); err != nil {
		return "", err
	}

	return u.String(), nil
}

// IsZero reports whether the UID is absent.
func (u UID) IsZero() bool {
	return u == ""
}

// IsAbsent reports whether the UID is absent.
func (u UID) IsAbsent() bool {
	return u == ""
}
