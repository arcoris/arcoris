/*
  Copyright 2026 The ARCORIS Authors

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
*/

package reduce

// Merger folds src into dst.
//
// Merge implementations call Merger after mapper execution has completed and do
// so from one goroutine. Merger may mutate only dst; src is a by-value partial.
// For floating-point reductions, callers should choose a merge mode with the
// expected rounding behavior for their algorithm because grouping can affect the
// final bits.
type Merger[T any] func(dst *T, src T)
