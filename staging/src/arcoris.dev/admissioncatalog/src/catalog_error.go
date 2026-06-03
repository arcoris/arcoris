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

package admissioncatalog

import "errors"

// ErrNilReasonRegistry identifies Catalog construction without a reason
// catalog.
//
// Catalog needs a ReasonRegistry so reason lookup, listing, and registration
// remain explicit owner-created operations instead of package-level globals.
var ErrNilReasonRegistry = errors.New("admissioncatalog: nil reason registry")

// ErrNilComponentRegistry identifies Catalog construction without a component
// catalog.
//
// Catalog requires a ComponentRegistry so component metadata remains part of the
// aggregate owner-created catalog, not an implicit global registry.
var ErrNilComponentRegistry = errors.New("admissioncatalog: nil component registry")

// ErrMismatchedKindRegistry identifies Catalog construction with a component
// registry backed by a different kind catalog.
//
// Catalog requires one shared KindRegistry reference so RegisterKind and
// RegisterComponent observe the same open-world kind catalog. Pointer identity
// is intentional here: two registries with equal descriptors still represent
// different owner-created catalogs and can diverge after construction.
var ErrMismatchedKindRegistry = errors.New("admissioncatalog: mismatched kind registry")
