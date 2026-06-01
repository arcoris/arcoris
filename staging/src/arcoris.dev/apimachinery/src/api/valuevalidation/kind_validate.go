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

package valuevalidation

import (
	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
)

// addKindMismatch records a concrete kind / descriptor type mismatch.
func (v *validator) addKindMismatch(
	path fieldpath.Path,
	actual value.Kind,
	expected value.Kind,
	code types.TypeCode,
) {
	v.addf(
		path,
		ErrKindMismatch,
		ErrorReasonKindMismatch,
		"value kind %s does not match descriptor %s; expected %s",
		actual,
		code,
		expected,
	)
}

// requireKind reports whether val has expected kind and records a diagnostic otherwise.
func (v *validator) requireKind(
	path fieldpath.Path,
	val value.Value,
	expected value.Kind,
	code types.TypeCode,
) bool {
	if val.Kind() == expected {
		return true
	}

	v.addKindMismatch(path, val.Kind(), expected, code)
	return false
}
