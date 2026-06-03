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

// nilCatalogPanic is the stable panic string for nil *Catalog method receivers.
//
// A nil Catalog is a wiring/configuration bug. Missing catalog entries are
// represented by lookup methods returning false.
const nilCatalogPanic = "admissioncatalog.Catalog: nil catalog"

// Catalog aggregates owner-created admission metadata registries.
//
// Catalog is a convenience boundary for reason, kind, and component catalog
// access. It is not global, not a runtime instance registry, does not store live
// Admitter values, and does not execute admission chains. Each method delegates
// to the underlying concurrency-safe registries.
type Catalog struct {
	// reasons stores stable reason descriptors for documentation, config
	// validation, and higher-level catalog checks.
	reasons *ReasonRegistry

	// kinds stores stable component kind descriptors.
	kinds *KindRegistry

	// components stores stable component descriptors validated against the same
	// kind registry held by Catalog.
	//
	// RegisterKind and RegisterComponent must observe one kind catalog. Passing
	// mismatched registries would split that invariant, so NewCatalog rejects a
	// ComponentRegistry backed by any other KindRegistry reference.
	components *ComponentRegistry
}
