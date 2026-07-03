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

// Package objectreflectorsink provides reusable implementations and
// combinators for runtime/objectreflector.Sink.
//
// The core objectreflector package owns the Sink contract and the
// source-to-sink protocol. This package contains concrete Sink helpers that can
// be used when wiring runtime components together without expanding the core
// reflector package.
//
// Fanout is a sequential Sink combinator. It preserves constructor order, stops
// on the first downstream error, and returns that error unchanged. Fanout does
// not provide atomicity across sinks, does not repair partial success, and does
// not retry failed sink operations. Callers should order sinks intentionally;
// for example, a read-model sink should usually run before a sink that exposes
// reconciliation work.
package objectreflectorsink
