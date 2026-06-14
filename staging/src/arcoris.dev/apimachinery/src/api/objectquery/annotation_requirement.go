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

package objectquery

// AnnotationRequirement is one canonical annotation selector requirement.
//
// The type is intentionally opaque. Callers must use the constructor functions
// so objectquery can validate annotation keys and values, canonicalize
// membership sets, and preserve deterministic predicate ordering.
type AnnotationRequirement struct {
	// req stores the shared metadata requirement representation.
	req metadataRequirement
}
