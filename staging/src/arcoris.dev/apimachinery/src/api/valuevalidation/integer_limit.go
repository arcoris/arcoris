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

package valuevalidation

// integerValue is the normalized integer domain used by valuevalidation.
//
// api/value.Integer carries the full int64 union uint64 domain. Descriptor
// validation narrows that union into either int64 or uint64 before it reaches
// this file, so range checks only need those two normalized forms.
type integerValue interface {
	~int64 | ~uint64
}

// integerBound represents one optional integer boundary.
type integerBound[T integerValue] struct {
	value T
	set   bool
}

// integerLimits groups optional inclusive integer boundaries.
type integerLimits[T integerValue] struct {
	lower integerBound[T]
	upper integerBound[T]
}
