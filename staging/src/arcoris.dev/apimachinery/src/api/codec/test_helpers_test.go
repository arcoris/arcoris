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

package codec

import (
	"errors"
	"reflect"
	"testing"
)

func validInfo() Info {
	return Info{
		Format:     FormatJSON,
		MediaTypes: []MediaType{MediaTypeJSON},
		Targets:    []Target{TargetValue},
	}
}

func requireNoError(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func requireErrorIs(t *testing.T, err error, target error) {
	t.Helper()

	if !errors.Is(err, target) {
		t.Fatalf("errors.Is(%v, %v) = false", err, target)
	}
}

func requireCodecError(t *testing.T, err error, path string, reason ErrorReason) {
	t.Helper()

	var codecErr *Error
	if !errors.As(err, &codecErr) {
		t.Fatalf("error type = %T; want *Error", err)
	}
	if codecErr.Path != path {
		t.Fatalf("Error.Path = %q; want %q", codecErr.Path, path)
	}
	if codecErr.Reason != reason {
		t.Fatalf("Error.Reason = %q; want %q", codecErr.Reason, reason)
	}
	if codecErr.Detail == "" {
		t.Fatalf("Error.Detail is empty")
	}
}

func requireMediaTypes(t *testing.T, got []MediaType, want ...MediaType) {
	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("media types = %#v; want %#v", got, want)
	}
}

func requireTargets(t *testing.T, got []Target, want ...Target) {
	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("targets = %#v; want %#v", got, want)
	}
}
