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

// validateUint64 checks TypeUint64 bounds, enum uniqueness, and enum membership.
func validateUint64(t Type, path string) error {
	if invalidRange(t.uint64.min, t.uint64.max) {
		return typeError(path+".range", ErrInvalidType)
	}
	if hasDuplicates(t.uint64.enum) || enumBelowMin(t.uint64.enum, t.uint64.min) || enumAboveMax(t.uint64.enum, t.uint64.max) {
		return typeError(path+".enum", ErrInvalidType)
	}
	return nil
}
