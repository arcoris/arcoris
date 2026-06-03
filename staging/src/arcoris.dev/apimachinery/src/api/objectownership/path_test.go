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

package objectownership

import (
	"testing"

	"arcoris.dev/apimachinery/api/fieldpath"
)

func TestPathString(t *testing.T) {
	if Path("$.image").String() != "$.image" {
		t.Fatalf("Path.String() = %q", Path("$.image").String())
	}
}

func TestPathParse(t *testing.T) {
	got, err := Path("$.image").Parse()
	requireNoError(t, err)

	if got.String() != "$.image" {
		t.Fatalf("Parse() = %q; want $.image", got)
	}
}

func TestPathParseRejectsEmpty(t *testing.T) {
	_, err := Path("").Parse()

	requireErrorIs(t, err, ErrInvalidPath)
}

func TestPathParseRejectsEmptyReportsPathDiagnostic(t *testing.T) {
	_, err := Path("").Parse()

	requireObjectOwnershipError(t, err, "path", ErrorReasonInvalidPath)
}

func TestPathParseRejectsInvalidPathReportsPathDiagnostic(t *testing.T) {
	_, err := Path("image").Parse()

	requireErrorIs(t, err, ErrInvalidPath)
	requireErrorIs(t, err, fieldpath.ErrInvalidPath)
	requireObjectOwnershipError(t, err, "path", ErrorReasonInvalidPath)
}

func TestPathParseRejectsNonCanonicalPath(t *testing.T) {
	_, err := Path(`$."image"`).Parse()

	requireErrorIs(t, err, ErrInvalidPath)
	requireObjectOwnershipError(t, err, "path", ErrorReasonInvalidPath)
}
