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

// Builder incrementally constructs semantic paths during hot traversal.
//
// Builder is mutable and not concurrency-safe. It should be kept local to a
// traversal and converted to an immutable Path with Path when a stable value is
// needed. Read-like methods are nil-safe; mutation methods require a non-nil
// builder and panic on nil receivers as programmer errors.
type Builder struct {
	elements []Element
}

// NewBuilder returns an empty path builder.
func NewBuilder() Builder {
	return Builder{}
}

// Len returns the number of currently pushed elements.
func (b *Builder) Len() int {
	if b == nil {
		return 0
	}

	return len(b.elements)
}

// Reset clears the builder while keeping its allocated storage for reuse.
//
// Reset requires a non-nil builder.
func (b *Builder) Reset() {
	b.elements = b.elements[:0]
}
