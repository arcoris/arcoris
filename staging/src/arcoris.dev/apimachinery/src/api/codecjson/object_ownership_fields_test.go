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

import (
	"testing"

	"arcoris.dev/apimachinery/api/apidocument"
)

func TestAllowOwnershipStateField(t *testing.T) {
	if !allowOwnershipStateField(apidocument.OwnershipFieldDesired.String()) || !allowOwnershipStateField(apidocument.OwnershipFieldMetadata.String()) {
		t.Fatalf("state field allow-list rejected known fields")
	}
	if allowOwnershipStateField(apidocument.OwnershipFieldEntries.String()) {
		t.Fatalf("state field allow-list accepted surface field")
	}
}

func TestAllowOwnershipSurfaceField(t *testing.T) {
	if !allowOwnershipSurfaceField(apidocument.OwnershipFieldEntries.String()) {
		t.Fatalf("surface field allow-list rejected entries")
	}
	if allowOwnershipSurfaceField(apidocument.OwnershipFieldOwner.String()) {
		t.Fatalf("surface field allow-list accepted entry field")
	}
}

func TestAllowOwnershipEntryField(t *testing.T) {
	if !allowOwnershipEntryField(apidocument.OwnershipFieldOwner.String()) || !allowOwnershipEntryField(apidocument.OwnershipFieldFields.String()) {
		t.Fatalf("entry field allow-list rejected known fields")
	}
	if allowOwnershipEntryField(apidocument.OwnershipFieldDesired.String()) {
		t.Fatalf("entry field allow-list accepted document field")
	}
}
