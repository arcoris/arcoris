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

// Package codecregistry provides owner-created registries for configured API
// codec implementations.
//
// The package indexes already-configured api/codec BaseCodec implementations by
// media type and groups them by format. MediaType is the unique lookup key.
// Format is a non-unique grouping attribute for families such as JSON, YAML, or
// CBOR, where one family may expose multiple media types. The registry validates
// codec metadata, checks declared targets against implemented byte and
// streaming capabilities, and exposes typed lookup helpers for codec
// capabilities.
//
// Registry construction accepts normalizable codec.Info metadata, stores a
// normalized detached metadata snapshot, and rejects invalid or non-normalizable
// metadata. Lookup and listing methods always operate on that normalized
// snapshot rather than calling codec.Info again.
//
// Registries are immutable after construction and do not use global mutable
// state, init-time registration, or default codec bundles. Callers decide which
// already-configured concrete implementations to register.
//
// A codec's Info.Targets must match its implemented capability interfaces.
// Declaring TargetValue requires ValueCodec or ValueStreamCodec. Implementing
// ValueCodec or ValueStreamCodec requires TargetValue. The same rule applies to
// TargetObject with ObjectCodec/ObjectStreamCodec and TargetObjectOwnership
// with ObjectOwnershipCodec/ObjectOwnershipStreamCodec.
//
// The package does not configure codecs, resolve options, interpret profiles,
// parse protocol metadata, negotiate media preferences, create codec instances,
// implement codecs, validate API values, apply objects, access storage, run
// admission, perform resource catalog lookup, or execute runtime/server
// lifecycle behavior.
package codecregistry
