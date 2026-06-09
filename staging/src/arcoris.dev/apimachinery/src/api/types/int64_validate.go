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

// validateInt64 checks DescriptorInt64 bounds, enum uniqueness, and enum membership.
func validateInt64(desc Descriptor, path string) error {
	if err := validateRangeRule(path, "int64", desc.int64.min, desc.int64.max); err != nil {
		return err
	}

	return validateEnumRules(path, "int64", desc.int64.enum, desc.int64.min, desc.int64.max)
}
