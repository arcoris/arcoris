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

// Package typeref provides shared DescriptorRef traversal guards for descriptor-aware
// API packages.
//
// The package is intentionally internal and diagnostic-neutral. Callers keep
// their own public sentinel and reason models, while this package owns the
// common resolver, recursion, and depth-limit mechanics.
//
// Resolver resolves one DescriptorRef edge at a time. Callers explicitly enter and
// leave the returned reference name while descending, which keeps comparison,
// validation, and field-set extraction free to preserve their own control flow
// and semantic paths. ResolveFinal is available for call sites that only need to
// inspect the non-reference descriptor behind a chain.
package typeref
