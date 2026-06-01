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

// Package diagnostic contains the shared record model for API diagnostics.
//
// Public API packages keep their own Error structs, sentinel errors, and
// reason types. This internal package centralizes the repeated record storage,
// construction, formatting, and unwrap behavior while keeping domain-specific
// error types in their owning packages.
package diagnostic
