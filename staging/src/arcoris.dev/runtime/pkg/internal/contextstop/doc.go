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

// Package contextstop preserves context stop causes for package-owned error
// classifiers.
//
// # Package scope
//
// contextstop owns only the narrow mechanics for retaining the most useful
// context cancellation cause. It does not decide whether a stop is a timeout,
// interruption, retry stop, task stop, or caller-owned error.
//
// # Relationship to adjacent packages
//
// wait and retry use this package to avoid duplicating context.Cause fallback
// behavior while keeping their own error classification contracts. Higher-level
// packages must not use contextstop as a generic context-error policy package.
//
// # File ownership
//
// cause.go owns context-cause preservation.
//
// # Dependency policy
//
// Production code depends only on the Go standard library.
package contextstop
