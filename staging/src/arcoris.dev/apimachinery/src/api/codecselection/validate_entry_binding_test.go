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
	"testing"

	"arcoris.dev/apimachinery/api/codec"
	"arcoris.dev/apimachinery/api/codecregistry"
)

func TestValidateEntryBindingAtRejectsMediaTypeMismatch(t *testing.T) {
	registry := testRegistry(
		t,
		testFullByteRegistration("json.public", codec.MediaTypeYAML),
	)
	entry, ok := registry.LookupID(codecregistry.MustEntryID("json.public"))
	if !ok {
		t.Fatalf("json.public entry missing")
	}

	err := validateEntryBindingAt(
		"codecselection.decodeBindings[0]",
		entry,
		testContentType(codec.MediaTypeJSON),
		codec.TargetObject,
		TransportBytes,
	)

	requireErrorIs(t, err, ErrEntryMediaTypeMismatch)
	requireSelectionError(t, err, "codecselection.decodeBindings[0].contentType", ErrorReasonEntryMediaTypeMismatch)
}

func TestValidateEntryBindingAtRejectsTargetMismatch(t *testing.T) {
	registry := testRegistry(
		t,
		testValueByteRegistration("json.public", codec.MediaTypeJSON),
	)
	entry, ok := registry.LookupID(codecregistry.MustEntryID("json.public"))
	if !ok {
		t.Fatalf("json.public entry missing")
	}

	err := validateEntryBindingAt(
		"codecselection.decodeBindings[0]",
		entry,
		testContentType(codec.MediaTypeJSON),
		codec.TargetObject,
		TransportBytes,
	)

	requireErrorIs(t, err, ErrEntryTargetMismatch)
	requireSelectionError(t, err, "codecselection.decodeBindings[0].target", ErrorReasonEntryTargetMismatch)
}

func TestValidateEntryBindingAtRejectsCapabilityMismatch(t *testing.T) {
	registry := testRegistry(
		t,
		testValueStreamRegistration("json.public", codec.MediaTypeJSON),
	)
	entry, ok := registry.LookupID(codecregistry.MustEntryID("json.public"))
	if !ok {
		t.Fatalf("json.public entry missing")
	}

	err := validateEntryBindingAt(
		"codecselection.decodeBindings[0]",
		entry,
		testContentType(codec.MediaTypeJSON),
		codec.TargetValue,
		TransportBytes,
	)

	requireErrorIs(t, err, ErrEntryCapabilityMismatch)
	requireSelectionError(t, err, "codecselection.decodeBindings[0].transport", ErrorReasonEntryCapabilityMismatch)
}
