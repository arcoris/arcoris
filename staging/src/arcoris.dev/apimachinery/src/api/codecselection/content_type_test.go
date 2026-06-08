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
)

func TestContentTypeNormalizesMediaTypeAndParameters(t *testing.T) {
	contentType, err := NewContentType(
		" Application/JSON ",
		Parameter{name: "Profile", value: "canonical"},
		Parameter{name: "version", value: "V1"},
	)
	requireNoError(t, err)

	if got := contentType.MediaType(); got != codec.MediaTypeJSON {
		t.Fatalf("media type = %q; want %q", got, codec.MediaTypeJSON)
	}

	items := contentType.Parameters().Items()
	if len(items) != 2 {
		t.Fatalf("parameters length = %d; want 2", len(items))
	}
	if items[0].Name() != "profile" || items[0].Value() != "canonical" {
		t.Fatalf("parameter[0] = %s=%s; want profile=canonical", items[0].Name(), items[0].Value())
	}
	if items[1].Name() != "version" || items[1].Value() != "V1" {
		t.Fatalf("parameter[1] = %s=%s; want version=V1", items[1].Name(), items[1].Value())
	}
	if got := contentType.String(); got != "application/json;profile=canonical;version=V1" {
		t.Fatalf("String() = %q; want normalized representation", got)
	}
}

func TestContentTypeRejectsDuplicateParameters(t *testing.T) {
	_, err := NewContentType(
		codec.MediaTypeJSON,
		Parameter{name: "profile", value: "canonical"},
		Parameter{name: "PROFILE", value: "storage"},
	)

	requireErrorIs(t, err, ErrInvalidParameters)
	requireSelectionError(t, err, "codecselection.contentType.parameters[1]", ErrorReasonInvalidParameters)
}

func TestContentTypeParametersReturnsDetachedData(t *testing.T) {
	contentType := MustContentType(
		codec.MediaTypeJSON,
		MustParameter("profile", "canonical"),
	)

	parameters := contentType.Parameters()
	items := parameters.Items()
	items[0] = MustParameter("profile", "storage")

	if got := contentType.Parameters().Items()[0].Value(); got != "canonical" {
		t.Fatalf("detached content type parameter mutation changed source: %q", got)
	}
}
