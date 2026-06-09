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

// Version identifies an API version.
//
// Valid versions are "vN", "vNalphaM", or "vNbetaM", where N and M are
// positive decimal integers without leading zeroes. Version is intentionally
// not SemVer and does not accept release-candidate tags, arbitrary labels, v0,
// uppercase forms, trimming, or normalization.
type Version string

// String returns the canonical version text without revalidating it.
func (v Version) String() string {
	return string(v)
}

// CanonicalText validates the version and returns its canonical text.
func (v Version) CanonicalText() (string, error) {
	if err := v.Validate(); err != nil {
		return "", err
	}

	return v.String(), nil
}

// IsZero reports whether the version is absent.
func (v Version) IsZero() bool {
	return v == ""
}

// IsAbsent reports whether the version is absent.
func (v Version) IsAbsent() bool {
	return v == ""
}
