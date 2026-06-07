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

// Package healthregistry provides in-process health check registration and
// resolution.
//
// The package implements arcoris.dev/health.CheckResolver with a mutable Builder
// and immutable Registry. Builder owns setup-time registration validation.
// Registry owns read-only target resolution after Build.
//
// Package healthregistry does not execute checks, aggregate reports, apply
// target policy, expose transports, run probes, schedule checks, collect
// metrics, log, trace, or make restart, admission, routing, or scheduler
// decisions. Those concerns belong to evaluator, probe, adapter, observability,
// and runtime-owner packages.
package healthregistry
