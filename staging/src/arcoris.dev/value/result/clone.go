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

package result

// Clone returns a Result with the successful value cloned by clone.
//
// Clone panics if clone is nil. The function is required even for Err results
// so clone configuration errors are detected consistently.
func (r Result[T]) Clone(clone func(T) T) Result[T] {
	if clone == nil {
		panic("result: nil clone function")
	}
	if r.err != nil {
		return Err[T](r.err)
	}
	return OK(clone(r.value))
}
