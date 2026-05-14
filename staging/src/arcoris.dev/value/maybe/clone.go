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

package maybe

// Clone returns a Maybe with the contained value cloned by clone.
//
// Clone panics if clone is nil. The function is required even for None values
// so clone configuration errors are detected consistently.
func (m Maybe[T]) Clone(clone func(T) T) Maybe[T] {
	if clone == nil {
		panic("maybe: nil clone function")
	}
	if !m.ok {
		return None[T]()
	}
	return Some(clone(m.value))
}
