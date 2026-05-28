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

// validateUint16 checks TypeUint16 bounds, enum uniqueness, and enum membership.
func validateUint16(t Type, path string) error {
	if invalidRange(t.uint16.min, t.uint16.max) {
		return typeError(path+".range", ErrInvalidType)
	}
	if hasDuplicates(t.uint16.enum) || enumBelowMin(t.uint16.enum, t.uint16.min) || enumAboveMax(t.uint16.enum, t.uint16.max) {
		return typeError(path+".enum", ErrInvalidType)
	}
	return nil
}
