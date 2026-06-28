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

package snapshot

import (
	"maps"
	"slices"
	"testing"
)

type mutableReadModel struct {
	Names  []string
	Attrs  map[string]string
	Nested *nestedReadModel
}

type nestedReadModel struct {
	Tags []string
}

func mutableReadModelValue(name, attr, tag string) mutableReadModel {
	return mutableReadModel{
		Names: []string{name},
		Attrs: map[string]string{
			"role": attr,
		},
		Nested: &nestedReadModel{
			Tags: []string{tag},
		},
	}
}

func cloneMutableReadModel(v mutableReadModel) mutableReadModel {
	out := mutableReadModel{
		Names: slices.Clone(v.Names),
		Attrs: maps.Clone(v.Attrs),
	}
	if v.Nested != nil {
		out.Nested = &nestedReadModel{
			Tags: slices.Clone(v.Nested.Tags),
		}
	}
	return out
}

func mutateMutableReadModel(v *mutableReadModel, name, attr, tag string) {
	v.Names[0] = name
	v.Attrs["role"] = attr
	v.Nested.Tags[0] = tag
}

func assertMutableReadModel(t *testing.T, got, want mutableReadModel) {
	t.Helper()

	if !slices.Equal(got.Names, want.Names) {
		t.Fatalf("Names = %#v, want %#v", got.Names, want.Names)
	}
	if !maps.Equal(got.Attrs, want.Attrs) {
		t.Fatalf("Attrs = %#v, want %#v", got.Attrs, want.Attrs)
	}
	if got.Nested == nil || want.Nested == nil {
		if got.Nested != want.Nested {
			t.Fatalf("Nested = %#v, want %#v", got.Nested, want.Nested)
		}
		return
	}
	if !slices.Equal(got.Nested.Tags, want.Nested.Tags) {
		t.Fatalf("Nested.Tags = %#v, want %#v", got.Nested.Tags, want.Nested.Tags)
	}
}

type clonePanicAfter[T any] struct {
	after int
	calls int
	clone func(T) T
}

func (c *clonePanicAfter[T]) Clone(v T) T {
	c.calls++
	if c.calls == c.after {
		panic("clone failed")
	}
	return c.clone(v)
}

func cloneStrings(v []string) []string {
	return slices.Clone(v)
}
