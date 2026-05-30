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

// ItemsCopy returns a shallow copy of the item slice.
//
// It does not clone item values. Nil items return nil; non-nil empty items
// return a distinct empty slice.
func (l List[T]) ItemsCopy() []T {
	return copyItems(l.Items)
}

// copyItems returns a shallow copy of items while preserving nil-vs-empty shape.
func copyItems[T any](items []T) []T {
	if items == nil {
		return nil
	}

	out := make([]T, len(items))
	copy(out, items)

	return out
}
