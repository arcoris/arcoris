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
	"errors"
	"testing"

	"arcoris.dev/apimachinery/api/codec"
	"arcoris.dev/apimachinery/api/codecregistry"
	"arcoris.dev/apimachinery/encoding/json/jsonconfig"
)

func TestRegistrationBuildsExplicitRegistryEntry(t *testing.T) {
	registration, err := Registration(codecregistry.MustEntryID("json.public"), jsonconfig.Default())
	if err != nil {
		t.Fatalf("Registration() error = %v", err)
	}

	registry, err := codecregistry.New(registration)
	if err != nil {
		t.Fatalf("codecregistry.New() error = %v", err)
	}

	entry, ok := registry.LookupID(codecregistry.MustEntryID("json.public"))
	if !ok {
		t.Fatalf("registry missing JSON entry")
	}
	if entry.Info().Format != codec.FormatJSON {
		t.Fatalf("entry format = %q; want %q", entry.Info().Format, codec.FormatJSON)
	}
	if got := entry.Codec(); got != registration.Codec() {
		t.Fatalf("entry codec = %v; want registration codec %v", got, registration.Codec())
	}
}

func TestRegistrationPreservesInvalidConfigError(t *testing.T) {
	config := jsonconfig.Default()
	config.Decode.Limits.MaxDepth = -1

	registration, err := Registration(codecregistry.MustEntryID("json.public"), config)
	if err == nil {
		t.Fatalf("Registration() error = nil")
	}
	if !registration.IsZero() {
		t.Fatalf("Registration() returned non-zero registration on error: %#v", registration)
	}
	if !errors.Is(err, jsonconfig.ErrInvalidConfig) {
		t.Fatalf("Registration() error = %v; want ErrInvalidConfig", err)
	}
}

func TestRegistrationLeavesEntryIDValidationToRegistry(t *testing.T) {
	registration, err := Registration(codecregistry.EntryID("INVALID ID"), jsonconfig.Default())
	if err != nil {
		t.Fatalf("Registration() error = %v", err)
	}

	_, err = codecregistry.New(registration)
	if err == nil {
		t.Fatalf("codecregistry.New() error = nil")
	}
	if !errors.Is(err, codecregistry.ErrInvalidEntryID) {
		t.Fatalf("codecregistry.New() error = %v; want ErrInvalidEntryID", err)
	}
}
