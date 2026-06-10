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

// Namespace identifies the namespace portion of object metadata.
//
// The empty namespace means namespace absence. It does not mean "default" and
// this package never applies namespace defaulting.
type Namespace string

// String returns the raw namespace text without validating it.
func (n Namespace) String() string {
	return string(n)
}

// CanonicalText validates the namespace and returns its canonical text.
//
// The empty namespace is valid and returns an empty string because namespace
// absence is part of the object-name domain.
func (n Namespace) CanonicalText() (string, error) {
	if err := n.ValidateLexical(); err != nil {
		return "", err
	}

	return n.String(), nil
}

// IsZero reports whether the namespace is absent.
func (n Namespace) IsZero() bool {
	return n == ""
}

// IsAbsent reports whether no namespace is present.
func (n Namespace) IsAbsent() bool {
	return n == ""
}
