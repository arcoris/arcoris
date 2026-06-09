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

type exactNumericDescriptorCase[T comparable] struct {
	descriptor  Descriptor
	code        DescriptorKind
	min         func(Descriptor) (T, bool)
	max         func(Descriptor) (T, bool)
	enum        func(Descriptor) []T
	wrong       func(Descriptor) bool
	wantMin     T
	wantMax     T
	wantFirst   T
	replaceWith T
}

func requireExactNumericDescriptor[T comparable](t *testing.T, tt exactNumericDescriptorCase[T]) {
	t.Helper()

	requireCode(t, tt.descriptor, tt.code)
	requireNullable(t, tt.descriptor, true)

	min, ok := tt.min(tt.descriptor)
	requireEqual(t, ok, true)
	requireEqual(t, min, tt.wantMin)

	max, ok := tt.max(tt.descriptor)
	requireEqual(t, ok, true)
	requireEqual(t, max, tt.wantMax)

	enum := tt.enum(tt.descriptor)
	enum[0] = tt.replaceWith
	requireEqual(t, tt.enum(tt.descriptor)[0], tt.wantFirst)

	requireEqual(t, tt.wrong(tt.descriptor), false)
	requireValidDescriptor(t, tt.descriptor, nil)
}

type enumPayloadCloneCase[T comparable, P any] struct {
	payload     P
	clone       func(P) P
	empty       func(P) bool
	enum        func(P) []T
	setEnum     func(*P, []T)
	zero        P
	wantFirst   T
	replaceWith T
}

func requireEnumPayloadCloneAndEmpty[T comparable, P any](t *testing.T, tt enumPayloadCloneCase[T, P]) {
	t.Helper()

	cloned := tt.clone(tt.payload)
	enum := tt.enum(cloned)
	enum[0] = tt.replaceWith
	tt.setEnum(&cloned, enum)

	requireEqual(t, tt.enum(tt.payload)[0], tt.wantFirst)
	requireEqual(t, tt.empty(tt.zero), true)
	requireEqual(t, tt.empty(tt.payload), false)
}
