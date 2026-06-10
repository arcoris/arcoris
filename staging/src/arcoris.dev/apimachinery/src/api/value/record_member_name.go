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

package value

// MemberName is one concrete record member name.
//
// Member names are payload keys only. They must be non-empty, but descriptor
// field-name grammar, map-key constraints, and resource-specific naming rules
// belong to descriptor-aware validation outside this package.
type MemberName string

// NewMemberName validates name as a concrete record member name.
func NewMemberName(name string) (MemberName, error) {
	memberName := MemberName(name)
	if err := memberName.ValidateLexical(); err != nil {
		return "", err
	}

	return memberName, nil
}

// MustMemberName validates name or panics.
//
// It is intended for tests and static fixtures where an empty member name is a
// programmer error. Runtime construction paths should use NewMemberName.
func MustMemberName(name string) MemberName {
	memberName, err := NewMemberName(name)
	if err != nil {
		panic(err)
	}

	return memberName
}

// IsZero reports whether n is the absent member-name value.
func (n MemberName) IsZero() bool {
	return n == ""
}

// String returns the member name text.
func (n MemberName) String() string {
	return string(n)
}

// ValidateLexical checks the base value-layer member-name invariant.
func (n MemberName) ValidateLexical() error {
	if n.IsZero() {
		return newError(
			pathMemberName,
			ErrEmptyMemberName,
			ErrorReasonEmptyMemberName,
			"record member name is empty",
		)
	}

	return nil
}
