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

package objectsurface

// Kind identifies a stable object surface.
type Kind string

const (
	// Desired is declarative user or manager intent.
	Desired Kind = "desired"

	// Observed is runtime, controller, or agent truth.
	Observed Kind = "observed"

	// MetadataLabels is the ObjectMeta labels map surface.
	MetadataLabels Kind = "metadata.labels"

	// MetadataAnnotations is the ObjectMeta annotations map surface.
	MetadataAnnotations Kind = "metadata.annotations"

	// MetadataFinalizers is reserved for lifecycle finalizer semantics.
	MetadataFinalizers Kind = "metadata.finalizers"

	// MetadataOwnerReferences is reserved for ownership-reference governance.
	MetadataOwnerReferences Kind = "metadata.ownerReferences"
)

// String returns stable surface text.
func (k Kind) String() string {
	return string(k)
}

// IsValid reports whether k is a known object surface.
func (k Kind) IsValid() bool {
	switch k {
	case Desired,
		Observed,
		MetadataLabels,
		MetadataAnnotations,
		MetadataFinalizers,
		MetadataOwnerReferences:
		return true
	default:
		return false
	}
}

// IsOwnable reports whether k is currently modeled by object ownership.
func (k Kind) IsOwnable() bool {
	switch k {
	case Desired,
		Observed,
		MetadataLabels,
		MetadataAnnotations:
		return true
	default:
		return false
	}
}

// IsMetadata reports whether k belongs to ObjectMeta.
func (k Kind) IsMetadata() bool {
	switch k {
	case MetadataLabels,
		MetadataAnnotations,
		MetadataFinalizers,
		MetadataOwnerReferences:
		return true
	default:
		return false
	}
}
