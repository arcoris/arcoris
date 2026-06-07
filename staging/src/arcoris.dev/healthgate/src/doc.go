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

// Package healthgate provides mutable owner-published health check state.
//
// A Gate is a concurrency-safe arcoris.dev/health.Checker that stores the latest
// arcoris.dev/health.Result published by its owner. Gates are useful when a
// component already knows its startup, readiness, overload, drain, dependency,
// or fatal state and wants evaluators to read that state cheaply.
//
// Package healthgate does not execute work, poll dependencies, create timers,
// expose endpoints, emit metrics, log, trace, or decide restart, admission,
// routing, scheduler, or lifecycle behavior.
package healthgate
