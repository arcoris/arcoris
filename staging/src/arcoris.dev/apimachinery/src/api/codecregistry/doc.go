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

// Package codecregistry provides immutable owner-created candidate catalogs for
// configured API codec implementations.
//
// The package indexes already-configured api/codec BaseCodec implementations by
// caller-supplied EntryID and groups them by media type and format. EntryID is
// the unique identity of one configured codec candidate. MediaType is a
// non-unique grouping key for wire representation. Format is a non-unique
// grouping key for codec families such as JSON, YAML, or CBOR. Multiple entries
// may intentionally share the same MediaType and Format when they represent
// different configured runtime policies.
//
// The registry validates codec metadata, checks declared targets against
// implemented byte and streaming capabilities, stores normalized detached
// metadata snapshots, and exposes typed candidate-list helpers for codec
// capabilities.
//
// Registry construction accepts normalizable codec.Info metadata, stores a
// normalized detached metadata snapshot, and rejects invalid or non-normalizable
// metadata. Lookup and listing methods always operate on that normalized
// snapshot rather than calling codec.Info again.
//
// Registries are immutable after construction and do not use global mutable
// state, init-time registration, or default codec bundles. Callers decide which
// already-configured concrete implementations to register and assign each
// configured candidate a stable EntryID.
//
// A codec's Info.Targets must match its implemented capability interfaces.
// Declaring TargetValue requires ValueCodec or ValueStreamCodec. Implementing
// ValueCodec or ValueStreamCodec requires TargetValue. The same rule applies to
// TargetObject with ObjectCodec/ObjectStreamCodec and TargetObjectOwnership
// with ObjectOwnershipCodec/ObjectOwnershipStreamCodec.
//
// The package does not select an operational codec. It does not configure
// codecs, resolve options, interpret roles, profiles, parameters, MIME headers,
// HTTP Accept values, media preferences, or runtime policies. It does not create
// codec instances, implement codecs, validate API values, apply objects, access
// storage, run admission, perform resource catalog lookup, or execute
// runtime/server lifecycle behavior. Future selection/profile layers choose
// exactly one candidate from registry candidate sets.
package codecregistry
