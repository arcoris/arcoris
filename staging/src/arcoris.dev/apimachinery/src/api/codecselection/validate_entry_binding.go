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

package codecselection

import (
	"arcoris.dev/apimachinery/api/codec"
	"arcoris.dev/apimachinery/api/codecregistry"
)

// validateEntryBindingAt checks registry metadata and typed capability support.
func validateEntryBindingAt(
	path string,
	entry codecregistry.Entry,
	contentType ContentType,
	target codec.Target,
	transport Transport,
) error {
	info := entry.Info()
	if !info.SupportsMediaType(contentType.mediaType) {
		return errorfAt(
			bindingContentTypePath(path),
			ErrEntryMediaTypeMismatch,
			ErrorReasonEntryMediaTypeMismatch,
			"entry %q does not declare media type %q",
			entry.ID(),
			contentType.mediaType,
		)
	}
	if !info.Supports(target) {
		return errorfAt(
			bindingTargetPath(path),
			ErrEntryTargetMismatch,
			ErrorReasonEntryTargetMismatch,
			"entry %q does not declare target %q",
			entry.ID(),
			target,
		)
	}
	if !entrySupportsCapability(entry, target, transport) {
		return errorfAt(
			bindingTransportPath(path),
			ErrEntryCapabilityMismatch,
			ErrorReasonEntryCapabilityMismatch,
			"entry %q does not implement %s transport for target %q",
			entry.ID(),
			transport,
			target,
		)
	}

	return nil
}
