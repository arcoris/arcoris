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

package fieldownership

// MaxOwnerLength bounds owner identity text stored in ownership state.
//
// Owner names should be stable identities such as "user-cli", "terraform",
// "status-controller", "arcoris.dev/controller", or "user:anton". They should
// not contain raw request IDs, timestamps, secrets, or unbounded instance
// identifiers. This package enforces only safe syntax and length constraints;
// higher layers own identity policy.
const MaxOwnerLength = 128

// Owner identifies one field-ownership participant.
//
// Owner is not an admission role, authorization subject, RBAC principal,
// runtime component instance, request actor, storage key, or audit identity.
// Higher layers decide which request subject may act as which Owner.
type Owner struct {
	text string
}

// NewOwner validates text and returns a controlled owner identity.
func NewOwner(text string) (Owner, error) {
	owner := Owner{text: text}
	if err := validateNewOwner(owner); err != nil {
		return Owner{}, err
	}

	return owner, nil
}

// MustOwner returns a valid owner or panics for invalid static fixtures.
func MustOwner(text string) Owner {
	owner, err := NewOwner(text)
	if err != nil {
		panic(err)
	}

	return owner
}

// String returns o as plain text.
func (o Owner) String() string {
	return o.text
}

// IsZero reports whether o is the empty owner identity.
func (o Owner) IsZero() bool {
	return o.text == ""
}

// Compare orders owner identities lexicographically by stable text.
func (o Owner) Compare(other Owner) int {
	switch {
	case o.text < other.text:
		return -1
	case o.text > other.text:
		return 1
	default:
		return 0
	}
}
