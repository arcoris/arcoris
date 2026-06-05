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

// Package lexical provides internal ASCII and identifier validation helpers
// used by ARCORIS API descriptor packages.
//
// The package is internal on purpose. It does not define public API concepts,
// domain-owned names, or exported descriptor types. Callers must wrap lexical
// failures in their own domain-specific errors so public packages continue to
// expose their own sentinels, reasons, paths, and diagnostic wording.
//
// The package intentionally uses byte-level ASCII checks instead of Unicode
// predicates. API descriptor identifiers are protocol tokens, not user-facing
// display names. Validation therefore performs no trimming, no normalization,
// no Unicode identifier expansion, no DNS lookup behavior, and no regexp-based
// matching dependency.
package lexical
