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

package meta

// Validate checks metadata-level apiVersion/kind consistency.
func (m TypeMeta) Validate() error {
	if m.IsZero() {
		return nil
	}
	if m.APIVersion.IsZero() {
		return invalid(
			"typeMeta.apiVersion",
			ErrInvalidTypeMeta,
			ErrorReasonEmptyValue,
			"apiVersion is required when kind is set",
		)
	}
	if m.Kind.IsZero() {
		return invalid(
			"typeMeta.kind",
			ErrInvalidTypeMeta,
			ErrorReasonEmptyValue,
			"kind is required when apiVersion is set",
		)
	}
	if err := m.APIVersion.Validate(); err != nil {
		return nested("typeMeta.apiVersion", ErrInvalidTypeMeta, err)
	}
	if err := m.Kind.Validate(); err != nil {
		return nested("typeMeta.kind", ErrInvalidTypeMeta, err)
	}
	return nil
}
