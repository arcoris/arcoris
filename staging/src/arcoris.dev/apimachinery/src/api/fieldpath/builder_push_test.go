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

package fieldpath

import "testing"

func TestBuilderPushHelpers(t *testing.T) {
	selector := MustSelector(testSelectorEntry("type", StringLiteral("Ready")))
	builder := NewBuilder()

	builder.PushField(testField("metadata"))
	builder.PushKey(testKey("labels"))
	requireNoError(t, builder.PushIndex(0))
	requireNoError(t, builder.PushSelector(selector))

	requireEqual(t, builder.Path().String(), `$.metadata["labels"][0][{"type":"Ready"}]`)
}

func TestBuilderPushIndexRejectsNegativeIndex(t *testing.T) {
	builder := NewBuilder()

	err := builder.PushIndex(-1)

	requireErrorIs(t, err, ErrNegativeIndex)
}

func TestNilBuilderMutationPanics(t *testing.T) {
	var builder *Builder

	requirePanic(t, func() {
		builder.PushField(testField("spec"))
	})
}

func TestBuilderPop(t *testing.T) {
	builder := NewBuilder()
	builder.PushField(testField("spec"))

	requireEqual(t, builder.Pop(), true)
	requireEqual(t, builder.Pop(), false)
	requireEqual(t, builder.Path().String(), "$")
}

func TestNilBuilderPopReturnsFalse(t *testing.T) {
	var builder *Builder

	requireEqual(t, builder.Pop(), false)
}
