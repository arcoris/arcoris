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

// Package codecregistry provides owner-created registries for API codec
// implementations.
//
// The package indexes api/codec BaseCodec implementations by format and media
// type, validates codec metadata, checks declared targets against implemented
// byte and streaming capabilities, and exposes typed lookup helpers for codec
// capabilities.
//
// Registries are immutable after construction and do not use global mutable
// state, init-time registration, or default codec bundles. Concrete codec
// implementations live in packages such as api/codecjson, api/codecyaml, and
// api/codeccbor; callers decide which implementations to register.
//
// A codec's Info.Targets must match its implemented capability interfaces.
// Declaring TargetValue requires ValueCodec or ValueStreamCodec. Implementing
// ValueCodec or ValueStreamCodec requires TargetValue. The same rule applies to
// TargetObject with ObjectCodec/ObjectStreamCodec and TargetObjectOwnership
// with ObjectOwnershipCodec/ObjectOwnershipStreamCodec.
//
// The package does not implement codecs, parse HTTP headers, negotiate Accept
// values, validate API values, apply objects, access storage, run admission,
// perform resource catalog lookup, or execute runtime/server lifecycle behavior.
package codecregistry
