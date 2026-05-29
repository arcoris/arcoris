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

// Package identity defines strict API identity primitives used by ARCORIS API
// machinery.
//
// The package provides immutable value types for API groups, versions, kinds,
// resources, subresources, and their canonical combinations. These values are
// lexical identity tokens. They do not carry runtime object behavior, metadata,
// structural type descriptors, codecs, REST routing policy, discovery
// documents, version negotiation, pluralization policy, or registration state.
//
// The package validates identity syntax only. Higher layers can use identity
// values to describe resources, bind Go runtime types, generate routes, choose
// codecs, or publish discovery output, but those responsibilities stay outside
// this package. In particular, identity is not a runtime scheme, metadata
// model, resource definition system, REST mapper, discovery document model, or
// Kubernetes compatibility layer.
//
// Each identity type is split by responsibility:
//
//   - the base file defines the value type, String/Identifier helpers,
//     IsZero semantics, getters, and pure composition methods;
//   - the parse file owns strict canonical text parsing;
//   - the validate file owns lexical validation and structured diagnostics;
//   - the encoding file owns text and JSON scalar encoding;
//   - grammar contains separator constants, while grammar helpers contain the
//     shared join/split mechanics used by String and Parse methods.
//
// That split is intentional. When changing one identity, start with the base
// value file to understand the shape, then follow the same name through parse,
// validate, encoding, and tests. Shared helpers remove mechanical repetition,
// but public behavior remains explicit on every concrete identity type.
//
// Canonical ARCORIS grammar uses one separator per identity axis:
//
//   - "/" separates API group and API version.
//   - "#" separates type identity from kind.
//   - ":" separates version identity from resource collection.
//   - "/" separates a resource collection from an optional subresource.
//
// Examples:
//
//	GroupVersion:
//
//		v1
//		control.arcoris.dev/v1
//
//	GroupKind:
//
//		Pod
//		control.arcoris.dev#Worker
//
//	GroupResource:
//
//		pods
//		control.arcoris.dev:workers
//
//	GroupVersionKind:
//
//		v1#Pod
//		control.arcoris.dev/v1#Worker
//
//	GroupVersionResource:
//
//		v1:pods
//		control.arcoris.dev/v1:workers
//
//	ResourcePath:
//
//		pods
//		pods/status
//
//	GroupVersionResourcePath:
//
//		v1:pods
//		v1:pods/status
//		control.arcoris.dev/v1:workers/status
//
// Parsing and encoding are strict round trips:
//
//	gvrp, err := ParseGroupVersionResourcePath(
//		"control.arcoris.dev/v1:workers/status",
//	)
//	if err != nil {
//		return err
//	}
//
//	text, err := gvrp.MarshalText()
//	if err != nil {
//		return err
//	}
//
// A direct literal can be useful in trusted package-internal declarations, but
// it is still validated before encoding:
//
//	gvk := GroupVersionKind{
//		Group:   "control.arcoris.dev",
//		Version: "v1",
//		Kind:    "Worker",
//	}
//
//	if err := gvk.Validate(); err != nil {
//		return err
//	}
//
// This package intentionally rejects Kubernetes-style legacy forms such as
// dotted kind/resource spellings, comma-based kind diagnostics, URL-like
// resource paths, and object-field helper APIs. Runtime object kind mutation,
// scheme registration, codec handling, resource definitions, and discovery or
// negotiation policy require separate explicit packages.
package identity
