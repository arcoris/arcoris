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

// Package objectcontrollerwiring assembles common object controller runtime
// graphs from lower-level apimachinery primitives.
//
// The package does not replace the lower-level packages. It records safe wiring
// patterns so controller authors do not repeat ordering-sensitive setup by
// hand. SameObject wires reflected object changes to reconciliation requests for
// the same object key.
//
// SameObject intentionally installs the cache sink before the enqueue sink.
// That order makes reflected state visible in the read model before matching
// work is made visible to controller workers.
//
// This package is a composition layer only. It constructs components, but does
// not own their lifecycle, launch goroutines, shut down queues, retry failures,
// requeue work, write objects, emit telemetry, recover panics, or implement
// scheduling policy.
package objectcontrollerwiring
