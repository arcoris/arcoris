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

package codecregistry

import (
	"testing"

	"arcoris.dev/apimachinery/api/codec"
)

func TestRegister(t *testing.T) {
	c := newValueByteCodec(codec.FormatJSON, codec.MediaTypeJSON)
	registration := Register(MustEntryID("json.public"), c)

	if registration.ID() != MustEntryID("json.public") {
		t.Fatalf("ID() = %q", registration.ID())
	}
	if registration.Codec() != c {
		t.Fatalf("Codec() = %v; want original codec", registration.Codec())
	}
}

func TestRegistrationIsZero(t *testing.T) {
	var registration Registration
	if !registration.IsZero() {
		t.Fatalf("zero registration IsZero() = false")
	}

	registration = Register(MustEntryID("json.public"), nil)
	if registration.IsZero() {
		t.Fatalf("registration with ID IsZero() = true")
	}
}

func TestRegistrationAccessors(t *testing.T) {
	c := newValueByteCodec(codec.FormatJSON, codec.MediaTypeJSON)
	registration := testRegistration("json.public", c)

	if registration.ID().String() != "json.public" {
		t.Fatalf("ID().String() = %q", registration.ID().String())
	}
	if registration.Codec() != c {
		t.Fatalf("Codec() = %v; want original codec", registration.Codec())
	}
}
