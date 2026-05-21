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


// Package core defines shared contracts for internal ARCORIS range reductions.
//
// A reduction splits an in-memory index interval into ranges, computes
// range-local or worker-local partial results, and merges those partials into
// one value. Core owns only the data model and callback contracts needed by that
// pattern: Range, Options, Strategy, MergeMode, mapper callbacks, Merger, and
// Scratch.
//
// Execution details live outside this package:
//
//   - planner builds deterministic non-overlapping ranges;
//   - runner executes range-local or worker-local reductions;
//   - merge folds completed partial results;
//   - layout keeps cache-line layout helpers tied to arcoris.dev/atomicx.
//
// Core must not import those implementation packages. The dependency direction
// is intentional: implementation packages depend on stable contracts, avoiding
// import cycles and keeping planning, execution, and merging separately
// testable. This package is not a public MapReduce API and does not expose a
// user-facing measurement surface.
package core
