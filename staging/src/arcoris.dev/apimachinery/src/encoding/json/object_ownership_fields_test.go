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

package jsoncodec

import (
	"testing"

	"arcoris.dev/apimachinery/api/apidocument"
)

func TestAllowOwnershipStateField(t *testing.T) {
	fields := apidocument.Fields().Ownership()

	if !allowOwnershipStateField(fields.Desired().String()) ||
		!allowOwnershipStateField(fields.Observed().String()) ||
		!allowOwnershipStateField(fields.MetadataField().String()) {
		t.Fatalf("state field allow-list rejected known fields")
	}
	if allowOwnershipStateField(fields.Surface().Entries().String()) {
		t.Fatalf("state field allow-list accepted surface field")
	}
}

func TestAllowOwnershipMetadataField(t *testing.T) {
	fields := apidocument.Fields().Ownership()
	metadataFields := fields.Metadata()

	if !allowOwnershipMetadataField(metadataFields.Labels().String()) ||
		!allowOwnershipMetadataField(metadataFields.Annotations().String()) {
		t.Fatalf("metadata field allow-list rejected known fields")
	}
	if allowOwnershipMetadataField(fields.Desired().String()) {
		t.Fatalf("metadata field allow-list accepted state field")
	}
}

func TestAllowOwnershipSurfaceField(t *testing.T) {
	fields := apidocument.Fields().Ownership()

	if !allowOwnershipSurfaceField(fields.Surface().Entries().String()) {
		t.Fatalf("surface field allow-list rejected entries")
	}
	if allowOwnershipSurfaceField(fields.Entry().Owner().String()) {
		t.Fatalf("surface field allow-list accepted entry field")
	}
}

func TestAllowOwnershipEntryField(t *testing.T) {
	fields := apidocument.Fields().Ownership()
	entryFields := fields.Entry()

	if !allowOwnershipEntryField(entryFields.Owner().String()) || !allowOwnershipEntryField(entryFields.Fields().String()) {
		t.Fatalf("entry field allow-list rejected known fields")
	}
	if allowOwnershipEntryField(fields.Desired().String()) {
		t.Fatalf("entry field allow-list accepted state field")
	}
}
