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

package codecjson

import "arcoris.dev/apimachinery/api/apidocument"

// allowOwnershipStateField reports whether name belongs to the ownership root.
func allowOwnershipStateField(name string) bool {
	switch name {
	case apidocument.OwnershipFieldDesired.String(),
		apidocument.OwnershipFieldObserved.String(),
		apidocument.OwnershipFieldMetadata.String():
		return true
	default:
		return false
	}
}

// allowOwnershipMetadataField reports whether name belongs to metadata ownership.
func allowOwnershipMetadataField(name string) bool {
	switch name {
	case apidocument.OwnershipFieldLabels.String(),
		apidocument.OwnershipFieldAnnotations.String():
		return true
	default:
		return false
	}
}

// allowOwnershipSurfaceField reports whether name belongs to a surface object.
func allowOwnershipSurfaceField(name string) bool {
	return name == apidocument.OwnershipFieldEntries.String()
}

// allowOwnershipEntryField reports whether name belongs to one ownership entry.
func allowOwnershipEntryField(name string) bool {
	switch name {
	case apidocument.OwnershipFieldOwner.String(),
		apidocument.OwnershipFieldFields.String():
		return true
	default:
		return false
	}
}
