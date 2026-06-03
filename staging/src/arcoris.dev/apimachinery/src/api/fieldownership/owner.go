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
// runtime component instance, or request actor. It is local ownership identity
// used only inside State.
type Owner string

// String returns o as plain text.
func (o Owner) String() string {
	return string(o)
}

// IsZero reports whether o is the empty owner identity.
func (o Owner) IsZero() bool {
	return o == ""
}
