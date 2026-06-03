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

// validateIdentityCompatibility rejects attempts to apply one object to another.
//
// The function treats version separately, so this helper checks only
// version-independent API kind identity, object name, namespace, and optional
// UID pinning.
func validateIdentityCompatibility(live ValueObject, applied ValueObject) error {
	liveGVK := live.GroupVersionKind()
	appliedGVK := applied.GroupVersionKind()

	// Group and kind mismatch means the applied object targets a different
	// resource family, even when metadata name happens to match.
	if liveGVK.Group != appliedGVK.Group || liveGVK.Kind != appliedGVK.Kind {
		return errorfAt(
			pathObjectAppliedTypeMeta,
			ErrIdentityMismatch,
			ErrorReasonIdentityMismatch,
			"applied GVK %s does not match live GVK %s",
			appliedGVK,
			liveGVK,
		)
	}

	liveName := live.ObjectName()
	appliedName := applied.ObjectName()

	// Namespace is part of ObjectName, so this rejects both renames and
	// cross-namespace apply attempts.
	if liveName != appliedName {
		return errorfAt(
			pathObjectAppliedMetadata,
			ErrIdentityMismatch,
			ErrorReasonIdentityMismatch,
			"applied object name %s does not match live object name %s",
			appliedName,
			liveName,
		)
	}

	appliedUID := applied.ObjectMeta.UID
	// Empty applied UID is allowed because clients may omit server-owned UID.
	// Non-empty UID must match the live object incarnation.
	if !appliedUID.IsZero() && appliedUID != live.ObjectMeta.UID {
		return errorfAt(
			pathObjectAppliedMetadata,
			ErrIdentityMismatch,
			ErrorReasonIdentityMismatch,
			"applied UID %q does not match live UID %q",
			appliedUID,
			live.ObjectMeta.UID,
		)
	}

	return nil
}
