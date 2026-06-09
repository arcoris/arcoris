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

// Subresource identifies an optional API subresource segment.
//
// The empty subresource is valid and means "no subresource". Non-empty
// subresources use the same DNS-1123 single-label grammar as Resource. The
// package does not predefine names such as "status" or "scale"; it validates
// lexical identity only.
type Subresource string

// String returns the canonical subresource text without revalidating it.
func (s Subresource) String() string {
	return string(s)
}

// CanonicalText validates the subresource and returns its canonical text.
//
// The empty subresource is valid and returns an empty string because absence is
// part of the subresource domain.
func (s Subresource) CanonicalText() (string, error) {
	if err := s.Validate(); err != nil {
		return "", err
	}

	return s.String(), nil
}

// IsZero reports whether the subresource is absent.
func (s Subresource) IsZero() bool {
	return s == ""
}

// IsAbsent reports whether no subresource segment is present.
func (s Subresource) IsAbsent() bool {
	return s == ""
}
