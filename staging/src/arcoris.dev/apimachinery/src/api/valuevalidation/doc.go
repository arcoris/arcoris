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

// Package valuevalidation validates concrete ARCORIS API payload values against
// structural api/types descriptors.
//
// The package is descriptor-aware and path-aware, but side-effect free. It does
// not decode wire formats, normalize values, apply defaults, prune fields,
// compare values, merge values, manage ownership, validate API object metadata,
// access storage, or perform admission.
//
// Validation errors use api/fieldpath semantic paths for payload locations.
// Object descriptor fields are addressed with field elements, dynamic map
// entries with key elements, ordered list items with index elements, and
// ListMap items with selector elements.
//
// Callers are expected to validate and register descriptors before using them
// for value validation. This package does not run full descriptor validation on
// every call. It performs only local defensive checks needed to avoid panics and
// to report invalid descriptor use encountered during traversal.
//
// List diagnostics intentionally use physical indexes for atomic, set-like, and
// ordered lists because item-level error locations are useful even when a later
// field-set/apply layer treats the complete list as one semantic field.
package valuevalidation
