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

// validateInt16 checks DescriptorInt16 bounds, enum uniqueness, and enum membership.
func validateInt16(desc Descriptor, path string) error {
	if err := validateRangeRule(path, "int16", desc.int16.min, desc.int16.max); err != nil {
		return err
	}

	return validateEnumRules(path, "int16", desc.int16.enum, desc.int16.min, desc.int16.max)
}
