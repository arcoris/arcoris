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


package core

// MergeMode selects how completed partial results are combined.
//
// Merge modes operate on already-computed partials. They do not affect planning
// or mapper scheduling.
type MergeMode uint8

const (
	// MergeLinear merges partial results from left to right in slice order.
	MergeLinear MergeMode = iota

	// MergePairwise merges partial results in pairwise rounds, reusing the
	// partial slice as working storage.
	MergePairwise
)
