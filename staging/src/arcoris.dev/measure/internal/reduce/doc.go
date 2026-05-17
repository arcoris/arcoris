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

// Package reduce defines the shared domain contracts for ARCORIS measurement
// range reductions.
//
// A reduction splits an in-memory index interval into ranges, computes
// range-local partial results, and merges those partials into one value. Root
// reduce owns only the data model and callback contracts needed by that pattern:
// Range, Options, strategies, mapper callbacks, merger callbacks, and Scratch.
//
// Execution details are deliberately kept in focused subpackages:
//
//   - planner builds deterministic non-overlapping ranges;
//   - runner executes range-local reducers;
//   - merge combines worker-local partial results;
//   - layout keeps cache-line layout helpers tied to arcoris.dev/atomicx.
//
// Root reduce does not import those implementations. The dependency direction is
// intentional: implementation packages depend on stable domain contracts, which
// avoids import cycles and keeps planning, execution, and merging separately
// testable. This package is not a public MapReduce API; it is an internal
// primitive for measurement packages that already know how to map and merge
// their own partial-result types.
package reduce
