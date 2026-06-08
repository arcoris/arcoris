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

// Package deadlineadmission maps deadline start decisions to admission results.
//
// The package is an adapter over arcoris.dev/resilience/deadline. It does not
// perform time math of its own, choose fallback timeout policy, create child
// contexts, wait, retry, schedule work, log, trace, export metrics, or propagate
// deadlines over transports. Core deadline math remains in package deadline.
//
// Admission results carry no grant and require no release. Metadata is the
// original deadline.Decision, so callers can inspect local deadline diagnostics
// without turning deadline reasons into wire-format compatibility contracts.
package deadlineadmission
