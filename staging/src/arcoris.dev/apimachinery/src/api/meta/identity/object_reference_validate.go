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

// ValidateLexical checks the typed name-only object reference metadata.
func (r ObjectReference) ValidateLexical() error {
	if r.APIVersion.IsZero() {
		return invalid(
			"objectReference.apiVersion",
			r.String(),
			ErrInvalidObjectReference,
			ErrorReasonEmptyValue,
			"apiVersion is required",
		)
	}

	if err := r.APIVersion.Validate(); err != nil {
		return nested("objectReference.apiVersion", r.String(), ErrInvalidObjectReference, err)
	}

	if r.Kind.IsZero() {
		return invalid(
			"objectReference.kind",
			r.String(),
			ErrInvalidObjectReference,
			ErrorReasonEmptyValue,
			"kind is required",
		)
	}

	if err := r.Kind.Validate(); err != nil {
		return nested("objectReference.kind", r.String(), ErrInvalidObjectReference, err)
	}

	if err := r.ObjectName().ValidateLexical(); err != nil {
		return nested("objectReference.objectName", r.String(), ErrInvalidObjectReference, err)
	}

	return nil
}
