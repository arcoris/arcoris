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

package objectapply

import (
	"arcoris.dev/apimachinery/api/meta"
	"arcoris.dev/apimachinery/api/meta/stamp"
)

// validateMetadataPolicy rejects applied metadata outside object identity.
//
// Identity has already been checked by validateIdentityCompatibility. This
// helper is only about metadata update intent that objectapply v1 does not
// support.
func validateMetadataPolicy(applied ValueObject) error {
	return rejectAppliedNonIdentityMetadata(applied.ObjectMeta)
}

// rejectAppliedNonIdentityMetadata enforces the Desired-only v1 metadata policy.
//
// Name, namespace, and optional UID are identity fields and may be present after
// identity compatibility has accepted them. Every other metadata field would be
// an object metadata mutation, so it is rejected instead of silently ignored.
func rejectAppliedNonIdentityMetadata(m meta.ObjectMeta) error {
	switch {
	case !m.GenerateName.IsZero():
		return unsupportedMetadata("generateName")
	case !m.ResourceVersion.IsZero():
		return unsupportedMetadata("resourceVersion")
	case !m.Generation.IsZero():
		return unsupportedMetadata("generation")
	case !m.CreatedAt.IsZero():
		return unsupportedMetadata("createdAt")
	case !deletionIsZero(m.Deletion):
		return unsupportedMetadata("deletion")
	case !m.Labels.IsZero():
		return unsupportedMetadata("labels")
	case !m.Annotations.IsZero():
		return unsupportedMetadata("annotations")
	case !m.OwnerReferences.IsZero():
		return unsupportedMetadata("ownerReferences")
	case !m.Finalizers.IsZero():
		return unsupportedMetadata("finalizers")
	default:
		return nil
	}
}

// deletionIsZero treats nil and empty deletion metadata as absent.
//
// A non-nil zero Deletion is not meaningful update intent, so it should not be
// rejected as a metadata change.
func deletionIsZero(deletion *stamp.Deletion) bool {
	return deletion == nil || deletion.IsZero()
}

// unsupportedMetadata reports an attempted metadata field change.
//
// The field argument is human-readable diagnostic context, not a JSON pointer or
// fieldpath.Path.
func unsupportedMetadata(field string) error {
	return errorfAt(
		pathObjectAppliedMetadata,
		ErrUnsupportedMetadataChange,
		ErrorReasonUnsupportedMetadataChange,
		"metadata field %q cannot be applied by objectapply v1",
		field,
	)
}
