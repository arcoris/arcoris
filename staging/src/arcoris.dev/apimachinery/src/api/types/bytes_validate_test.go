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

import "testing"

func TestBytesValidateRejectsInvalidLengthRules(t *testing.T) {
	negativeMin := Bytes().Descriptor()
	negativeMin.bytes.minBytes = limit[int]{value: -1, set: true}

	tests := []Descriptor{
		negativeMin,
		Bytes().MinBytes(2).MaxBytes(1).Descriptor(),
	}
	for _, desc := range tests {
		requireErrorIs(t, ValidateLocal(desc), ErrInvalidDescriptor)
	}
}
