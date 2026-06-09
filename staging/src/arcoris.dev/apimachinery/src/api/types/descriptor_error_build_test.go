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

import (
	"errors"
	"testing"
)

func TestDescriptorErrorBuildWrapsClassifiedError(t *testing.T) {
	err := descriptorError("descriptor.field", ErrInvalidDescriptor)
	requireEqual(t, errors.Is(err, ErrInvalidDescriptor), true)
	requireEqual(t, err.Error(), "types: descriptor.field: invalid descriptor")

	withoutPath := descriptorError("", ErrInvalidDescriptor)
	requireEqual(t, withoutPath.Error(), "types: invalid descriptor")
}

func TestDescriptorErrorBuildFormatsReasonAndDetail(t *testing.T) {
	err := descriptorErrorf(
		"descriptor.range",
		ErrInvalidDescriptor,
		DescriptorErrorReasonInvalidRange,
		"minimum must be <= maximum",
	)

	requireEqual(t, errors.Is(err, ErrInvalidDescriptor), true)
	requireEqual(t, err.Error(), "types: descriptor.range: invalid descriptor: invalid_range: minimum must be <= maximum")
}
