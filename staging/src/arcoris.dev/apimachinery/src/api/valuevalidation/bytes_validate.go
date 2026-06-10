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

package valuevalidation

import (
	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
)

// validateBytes checks DescriptorBytes kind and length constraints.
func (v *validator) validateBytes(path fieldpath.Path, val value.Value, descriptor types.Descriptor) {
	if !v.requireKind(path, val, value.KindBytes, descriptor.Code()) {
		return
	}

	bytes, _ := val.AsBytes()
	bytesView, ok := descriptor.AsBytes()
	if !ok {
		v.add(
			path,
			ErrInvalidDescriptor,
			ErrorReasonInvalidDescriptor,
			"descriptor is not bytes",
		)
		return
	}

	v.validateBytesLength(path, len(bytes), bytesView)
}

// validateBytesLength checks byte-sequence length rules.
func (v *validator) validateBytesLength(path fieldpath.Path, length int, bytesView types.BytesView) {
	if minLength, ok := bytesView.MinBytes(); ok && length < minLength {
		v.addf(
			path,
			ErrLengthOutOfRange,
			ErrorReasonTooShort,
			"bytes length %d is below minimum %d",
			length,
			minLength,
		)
	}

	if maxLength, ok := bytesView.MaxBytes(); ok && length > maxLength {
		v.addf(
			path,
			ErrLengthOutOfRange,
			ErrorReasonTooLong,
			"bytes length %d is above maximum %d",
			length,
			maxLength,
		)
	}
}
