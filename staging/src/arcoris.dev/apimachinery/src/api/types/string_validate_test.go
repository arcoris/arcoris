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

func TestStringValidateRejectsInvalidRules(t *testing.T) {
	tests := []Descriptor{
		String().MinBytes(-1).Descriptor(),
		String().MaxBytes(-1).Descriptor(),
		String().MinBytes(2).MaxBytes(1).Descriptor(),
		String().MinRunes(-1).Descriptor(),
		String().MaxRunes(-1).Descriptor(),
		String().MinRunes(2).MaxRunes(1).Descriptor(),
		String().Pattern("[").Descriptor(),
		String().MinBytes(2).Enum("a").Descriptor(),
		String().MaxBytes(1).Enum("ab").Descriptor(),
		String().MinRunes(2).Enum("a").Descriptor(),
		String().MaxRunes(1).Enum("ab").Descriptor(),
		String().MaxBytes(1).Enum("é").Descriptor(),
		String().Pattern("^a+$").Enum("bbb").Descriptor(),
		String().Enum("a", "a").Descriptor(),
	}
	for _, desc := range tests {
		requireErrorIs(t, ValidateLocal(desc), ErrInvalidDescriptor)
	}
}

func TestStringValidateDistinguishesBytesAndRunes(t *testing.T) {
	requireErrorIs(t, ValidateLocal(String().MaxBytes(1).Enum("é").Descriptor()), ErrInvalidDescriptor)
	requireNoError(t, ValidateLocal(String().MaxRunes(1).Enum("é").Descriptor()))
}

func TestStringValidateReportsRuneEnumErrors(t *testing.T) {
	requireDescriptorError(
		t,
		ValidateLocal(String().MinRunes(2).Enum("a").Descriptor()),
		ErrInvalidDescriptor,
		"descriptor.enum[0]",
		DescriptorErrorReasonEnumBelowMinimum,
		"rune count",
	)
	requireDescriptorError(
		t,
		ValidateLocal(String().MaxRunes(1).Enum("ab").Descriptor()),
		ErrInvalidDescriptor,
		"descriptor.enum[0]",
		DescriptorErrorReasonEnumAboveMaximum,
		"rune count",
	)
}
