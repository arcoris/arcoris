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

package types

// objectPayload stores fixed object fields and unknown-field policy.
//
// Fields are stored as a slice, not a map, because descriptor order is useful
// for diagnostics, generated code, schema export, and stable review output.
// Duplicate names are rejected by ValidateType so builders stay cheap and can
// report errors with complete descriptor paths.
type objectPayload struct {
	// fields contains finalized object fields in declaration order.
	//
	// Public accessors clone this slice to preserve descriptor immutability.
	fields []FieldDescriptor
	// unknown records the structural policy for undeclared fields.
	//
	// This package records the policy only; it does not process concrete object
	// values.
	unknown UnknownFieldPolicy
}

// cloneObjectPayload detaches ordered object fields.
func cloneObjectPayload(p objectPayload) objectPayload {
	p.fields = cloneFields(p.fields)
	return p
}
