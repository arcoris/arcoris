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

package capacity

// Refusal classifies an accounting refusal.
//
// Success is represented by ok=true, either through raw accounting or through an
// owned reservation. Refusal values are diagnostics for valid attempts that did
// not reserve capacity.
type Refusal uint8

const (
	// RefusalNone means local accounting did not refuse the attempt.
	RefusalNone Refusal = iota

	// RefusalInsufficient means known capacity exists but available capacity is too small.
	RefusalInsufficient

	// RefusalDebt means the demanded resource is already overcommitted.
	RefusalDebt

	// RefusalUnknownResource means vector accounting has no limit for a demanded resource.
	RefusalUnknownResource
)

// IsValid reports whether r is a known refusal value.
func (r Refusal) IsValid() bool {
	switch r {
	case RefusalNone, RefusalInsufficient, RefusalDebt, RefusalUnknownResource:
		return true
	default:
		return false
	}
}

// Refused reports whether r describes an accounting refusal.
func (r Refusal) Refused() bool {
	return r != RefusalNone && r.IsValid()
}

// String returns a stable diagnostic string.
func (r Refusal) String() string {
	switch r {
	case RefusalNone:
		return "none"
	case RefusalInsufficient:
		return "insufficient"
	case RefusalDebt:
		return "debt"
	case RefusalUnknownResource:
		return "unknown_resource"
	default:
		return "invalid"
	}
}
