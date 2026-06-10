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

package meta

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"
)

func TestPageToken(t *testing.T) {
	requireNoError(t, PageToken("").ValidateLexical())

	token, err := ParsePageToken("page-1")
	requireNoError(t, err)
	if token.String() != "page-1" || token.IsZero() {
		t.Fatalf("token = %q zero=%v", token, token.IsZero())
	}

	optional, err := ParseOptionalPageToken("")
	requireNoError(t, err)
	if !optional.IsZero() {
		t.Fatalf("optional token = %q", optional)
	}

	_, err = ParsePageToken("")
	requireErrorIs(t, err, ErrInvalidPageToken)
	_, err = ParsePageToken("page 1")
	requireErrorIs(t, err, ErrInvalidPageToken)
	_, err = ParsePageToken("page/1")
	requireErrorIs(t, err, ErrInvalidPageToken)
	_, err = ParsePageToken("page\n1")
	requireErrorIs(t, err, ErrInvalidPageToken)
	_, err = ParsePageToken(strings.Repeat("x", maxPageTokenLength+1))
	requireErrorIs(t, err, ErrInvalidPageToken)

	var parsed PageToken
	requireNoError(t, parsed.UnmarshalText([]byte("page-1")))
	if parsed != "page-1" {
		t.Fatalf("parsed = %q", parsed)
	}

	data, err := json.Marshal(PageToken("page-1"))
	requireNoError(t, err)
	if string(data) != `"page-1"` {
		t.Fatalf("json = %s", data)
	}

	err = json.Unmarshal([]byte(`null`), &parsed)
	requireErrorIs(t, err, ErrInvalidJSON)

	err = json.Unmarshal([]byte(`123`), &parsed)
	requireErrorIs(t, err, ErrInvalidJSON)

	err = json.Unmarshal([]byte(`"bad token"`), &parsed)
	requireErrorIs(t, err, ErrInvalidPageToken)
}

func TestPageTokenValidateStructuredLengthError(t *testing.T) {
	err := PageToken(strings.Repeat("x", maxPageTokenLength+1)).ValidateLexical()
	requireErrorIs(t, err, ErrInvalidPageToken)

	var metaErr *Error
	if !errors.As(err, &metaErr) {
		t.Fatalf("errors.As(%T) = false", metaErr)
	}
	if metaErr.Path != "pageToken" {
		t.Fatalf("Path = %q", metaErr.Path)
	}
	if metaErr.Reason != ErrorReasonInvalidLength {
		t.Fatalf("Reason = %q", metaErr.Reason)
	}
	if metaErr.Detail == "" {
		t.Fatal("Detail is empty")
	}
}
