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

package object

import "arcoris.dev/apimachinery/api/meta"

// WithTypeMeta returns a copy of the list with replacement type metadata.
func (l List[T]) WithTypeMeta(typeMeta meta.TypeMeta) List[T] {
	l.TypeMeta = typeMeta.Clone()

	return l
}

// WithPageMeta returns a copy of the list with replacement page metadata.
func (l List[T]) WithPageMeta(pageMeta meta.PageMeta) List[T] {
	l.PageMeta = pageMeta.Clone()

	return l
}

// WithItems returns a copy of the list with replacement items.
//
// The item slice is shallow-copied and nil-vs-empty shape is preserved. Item
// values themselves are not deep-copied.
func (l List[T]) WithItems(items []T) List[T] {
	l.Items = copyItems(items)

	return l
}
