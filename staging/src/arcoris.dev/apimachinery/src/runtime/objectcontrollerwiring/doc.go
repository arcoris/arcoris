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

// Package objectcontrollerwiring assembles small object-controller runtime
// graphs from lower-level apimachinery primitives.
//
// The package is a narrow composition layer. It does not replace the component
// packages it wires together; it records the ordering-sensitive setup that
// should be shared by controllers with the same shape.
//
// SameObject wires one watched source to reconciliation requests for the same
// object key. MappedObject wires one watched source to reconciliation requests
// produced by caller-provided objectenqueue mappers; one source item or change
// may emit zero, one, or many target work items. MultiSource wires a primary
// watched input plus optional secondary watched inputs into one shared queue and
// one controller.
//
// SameObject intentionally installs the cache sink before the enqueue sink.
// That order makes reflected state visible in the read model before matching
// work is made visible to controller workers.
//
// MappedObject's cache is the watched source collection cache. A mapped
// reconciler receives a target Request together with a snapshot of that source
// collection. The package does not add target caches or implicit secondary read
// paths.
//
// MultiSource uses the primary input cache as the controller snapshot source.
// Secondary input caches remain separate and are not merged into
// objectreconciler.Snapshot. Reconcilers that need secondary or target state
// should receive those dependencies explicitly.
//
// RunSameObject, RunMappedObject, and RunMultiSource coordinate one already
// assembled graph. They start the reflector or reflectors and controller
// together, shut down the queue when any component exits, cancel sibling
// components, and wait for everything before returning. This is the graph-level
// producer/consumer boundary: reflectors produce queue items, the controller
// drains them, and this package connects producer completion to queue shutdown.
//
// This package is intentionally not a manager. It does not register multiple
// controllers, retry, restart, requeue, supervise policy, write objects, patch
// status, emit telemetry, recover panics, or implement scheduling policy.
package objectcontrollerwiring
