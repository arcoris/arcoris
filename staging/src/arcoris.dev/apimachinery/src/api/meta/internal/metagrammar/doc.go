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

// Package metagrammar contains shared grammar helpers for metadata tokens.
//
// The package is internal implementation infrastructure for api/meta
// subpackages. It deliberately does not define public metadata concepts such as
// names, labels, annotations, or finalizers. Domain packages wrap failures in
// their own public sentinel errors and structured diagnostics.
//
// Metadata identifiers are protocol tokens, not display names. Validation uses
// byte-level ASCII rules, does not trim input, and does not normalize case.
package metagrammar
