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

// Package objectstore defines contracts for authoritative committed API object
// state.
//
// An object store holds value-backed live object envelopes, object ownership
// documents, and store-local commit revisions. It is descriptor-agnostic:
// callers are expected to resolve resources, validate object envelopes, compute
// apply results, and decide admission before committing state here.
//
// The package deliberately does not validate objects against resource
// contracts, apply objects, compute field conflicts, run admission, encode or
// decode wire formats, expose HTTP or gRPC behavior, list objects, watch
// changes, or stamp object metadata resourceVersion/generation fields. Those
// responsibilities belong to higher API operation, codec, request-flow,
// server, and future watch layers.
//
// This package defines contracts only. The in-memory implementation lives in
// the sibling package arcoris.dev/apimachinery/api/objectmemorystore.
package objectstore
