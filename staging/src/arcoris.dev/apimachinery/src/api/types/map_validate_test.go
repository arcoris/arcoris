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

func TestMapValidateRejectsInvalidShapes(t *testing.T) {
	requireErrorIs(t, ValidateLocal(MapOf(DescriptorExpr(nil)).Descriptor()), ErrInvalidDescriptor)
	requireErrorIs(t, ValidateLocal(MapOf(String()).MinEntries(2).MaxEntries(1).Descriptor()), ErrInvalidDescriptor)

	invalidKey := MapOf(String()).Descriptor()
	invalidKey.mapType.key = nil
	requireErrorIs(t, ValidateLocal(invalidKey), ErrInvalidDescriptor)

	boolKey := MapOf(String()).Keys(Bool()).Descriptor()
	requireErrorIs(t, ValidateLocal(boolKey), ErrInvalidField)

	nullableKey := MapOf(String()).Keys(String().Nullable()).Descriptor()
	requireErrorIs(t, ValidateLocal(nullableKey), ErrInvalidField)

	missingValue := Descriptor{code: DescriptorMap}
	key := String().Descriptor()
	missingValue.mapType.key = &key
	requireErrorIs(t, ValidateLocal(missingValue), ErrInvalidDescriptor)
}

func TestMapValidateResolvedKeyReference(t *testing.T) {
	keyDef := Define("meta.arcoris.dev.LabelKey", String().MinBytes(1))
	desc := MapOf(String()).Keys(Ref("meta.arcoris.dev.LabelKey")).Descriptor()

	requireNoError(t, ValidateResolved(desc, resolverFunc(func(name TypeName) (Definition, bool) {
		if name == keyDef.Name() {
			return keyDef, true
		}

		return Definition{}, false
	})))
}

func TestMapValidateRejectsUnresolvedKeyReference(t *testing.T) {
	desc := MapOf(String()).Keys(Ref("meta.arcoris.dev.LabelKey")).Descriptor()

	requireErrorIs(t, ValidateResolved(desc, resolverFunc(func(TypeName) (Definition, bool) {
		return Definition{}, false
	})), ErrUnresolvedDescriptorReference)
}

func TestMapValidateRejectsKeyReferenceToNonStringDescriptor(t *testing.T) {
	keyDef := Define("meta.arcoris.dev.LabelKey", Bool())
	desc := MapOf(String()).Keys(Ref("meta.arcoris.dev.LabelKey")).Descriptor()

	requireErrorIs(t, ValidateResolved(desc, resolverFunc(func(name TypeName) (Definition, bool) {
		if name == keyDef.Name() {
			return keyDef, true
		}

		return Definition{}, false
	})), ErrInvalidField)
}
