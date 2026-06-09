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

// Package objectlifecycle coordinates local API object lifecycle transitions
// over api/objectstore.
//
// The package is the first descriptor-aware stateful layer above the lower API
// machinery primitives. It resolves resource contracts, validates value-backed
// objects with api/objectvalidation, delegates Desired apply semantics to
// api/objectapply, and commits already-computed state through api/objectstore.
//
// Layering is deliberately narrow:
// objectvalidation validates object shape against an already resolved resource
// contract; objectapply computes a pure Desired apply result; objectstore
// commits object, ownership, and store revision with optimistic concurrency;
// objectlifecycle composes those pieces into Create, Apply, Delete, and Get
// operations.
//
// The package does not decode wire formats, select codecs, serve HTTP or gRPC,
// run admission, authorize subjects, list objects, watch resources, compact
// tombstones, stamp resourceVersion/generation metadata, generate UIDs, run
// finalizers, default, convert, prune, reconcile controllers, export metrics,
// log, trace, or start background goroutines.
package objectlifecycle
