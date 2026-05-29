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

package resourcecatalog

// ErrorReason identifies the precise catalog invariant that failed.
//
// Reasons are intentionally separate from sentinel errors. Callers can use
// errors.Is for broad handling and inspect ErrorReason for stable diagnostics.
type ErrorReason string

const (
	// ErrorReasonInvalidDefinition means a resource.Definition failed resource
	// descriptor validation during registration.
	ErrorReasonInvalidDefinition ErrorReason = "invalid_definition"

	// ErrorReasonDuplicateResource means one batch contains the same
	// GroupResource identity more than once.
	ErrorReasonDuplicateResource ErrorReason = "duplicate_resource"

	// ErrorReasonDuplicateKind means one batch contains the same GroupKind
	// identity more than once.
	ErrorReasonDuplicateKind ErrorReason = "duplicate_kind"

	// ErrorReasonDuplicateVersionResource means one batch contains the same
	// GroupVersionResource identity more than once.
	ErrorReasonDuplicateVersionResource ErrorReason = "duplicate_version_resource"

	// ErrorReasonDuplicateVersionKind means one batch contains the same
	// GroupVersionKind identity more than once.
	ErrorReasonDuplicateVersionKind ErrorReason = "duplicate_version_kind"

	// ErrorReasonDefinitionExists means an incoming definition conflicts with an
	// identity already stored in the catalog.
	ErrorReasonDefinitionExists ErrorReason = "definition_exists"

	// ErrorReasonNilCatalog means a mutating operation was called on a nil
	// Catalog receiver.
	ErrorReasonNilCatalog ErrorReason = "nil_catalog"
)
