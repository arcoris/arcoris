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

package typeref

// FailureKind classifies TypeRef traversal failures for package-local mapping.
type FailureKind string

const (
	// FailureInvalidDescriptor reports a descriptor that is not a valid TypeRef
	// at a call site that requires one.
	FailureInvalidDescriptor FailureKind = "invalid_descriptor"

	// FailureUnresolvedRef reports a missing resolver or a resolver miss.
	FailureUnresolvedRef FailureKind = "unresolved_ref"

	// FailureReferenceCycle reports a recursive or too-deep reference chain.
	FailureReferenceCycle FailureKind = "reference_cycle"
)
