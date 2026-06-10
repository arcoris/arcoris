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

package listmapkey

import (
	"strings"
	"testing"

	"arcoris.dev/apimachinery/api/value"
)

func TestKeyKindMismatchDetailNamesKinds(t *testing.T) {
	err := keyKindMismatch(
		conditionPath(0).Field(testFieldName("type")),
		value.KindInteger,
		value.KindString,
	)

	var extractionError *Error
	if !asErrorInto(err, &extractionError) {
		t.Fatalf("error type = %T, want *Error", err)
	}

	if !strings.Contains(extractionError.Detail, "integer") {
		t.Fatalf("detail %q does not include actual kind", extractionError.Detail)
	}

	if !strings.Contains(extractionError.Detail, "string") {
		t.Fatalf("detail %q does not include expected kind", extractionError.Detail)
	}
}

func asErrorInto(err error, target **Error) bool {
	extractionError, ok := asError(err)
	if !ok {
		return false
	}

	*target = extractionError

	return true
}
