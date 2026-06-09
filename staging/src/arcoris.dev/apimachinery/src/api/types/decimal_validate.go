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

// validateDecimal checks DescriptorDecimal precision and scale descriptor rules.
//
// Decimal deliberately has no min/max or enum rules in this package. Future
// descriptor bounds must be defined with exact decimal semantics and without
// binary floating-point conversion.
//
// Scale without precision is allowed. It constrains fractional shape while
// leaving total significant digits unconstrained.
func validateDecimal(desc Descriptor, path string) error {
	if desc.decimal.precision.set && desc.decimal.precision.value <= 0 {
		return descriptorErrorf(
			path+".precision",
			ErrInvalidDescriptor,
			DescriptorErrorReasonInvalidPrecision,
			"decimal precision must be greater than zero, got %d",
			desc.decimal.precision.value,
		)
	}

	if desc.decimal.scale.set && desc.decimal.scale.value < 0 {
		return descriptorErrorf(
			path+".scale",
			ErrInvalidDescriptor,
			DescriptorErrorReasonInvalidScale,
			"decimal scale must be non-negative, got %d",
			desc.decimal.scale.value,
		)
	}

	if desc.decimal.precision.set &&
		desc.decimal.scale.set &&
		desc.decimal.scale.value > desc.decimal.precision.value {
		return descriptorErrorf(
			path+".scale",
			ErrInvalidDescriptor,
			DescriptorErrorReasonInvalidScale,
			"decimal scale must be <= precision, got scale=%d precision=%d",
			desc.decimal.scale.value,
			desc.decimal.precision.value,
		)
	}

	return nil
}
