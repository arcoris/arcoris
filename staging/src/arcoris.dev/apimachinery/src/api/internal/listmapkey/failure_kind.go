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

package listmapkey

// FailureKind classifies why ListMap key extraction failed.
//
// Callers use the classification to decide whether they can fall back to
// physical index diagnostics, or whether the descriptor/resolver state itself
// must be reported.
type FailureKind string

const (
	// FailureInvalidDescriptor reports malformed descriptor state that prevents
	// selector extraction.
	FailureInvalidDescriptor FailureKind = "invalid_descriptor"

	// FailureUnresolvedRef reports a DescriptorRef that cannot be resolved.
	FailureUnresolvedRef FailureKind = "unresolved_ref"

	// FailureReferenceCycle reports recursive or too-deep DescriptorRef traversal.
	FailureReferenceCycle FailureKind = "reference_cycle"

	// FailureItemKindMismatch reports a ListMap item that is not a concrete
	// object payload.
	FailureItemKindMismatch FailureKind = "item_kind_mismatch"

	// FailureMissingKey reports an item that does not carry a selector key
	// member required by the descriptor.
	FailureMissingKey FailureKind = "missing_key"

	// FailureNullKey reports a selector key member whose concrete value is null.
	FailureNullKey FailureKind = "null_key"

	// FailureKeyKindMismatch reports a selector key value with the wrong
	// concrete kind for its descriptor.
	FailureKeyKindMismatch FailureKind = "key_kind_mismatch"

	// FailureKeyIntegerRange reports an integer selector key that does not fit
	// the signedness required by the key descriptor.
	FailureKeyIntegerRange FailureKind = "key_integer_range"

	// FailureInvalidSelector reports selector construction failure after key
	// extraction has otherwise succeeded.
	FailureInvalidSelector FailureKind = "invalid_selector"
)
