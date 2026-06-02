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
	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/value"
)

// stringLiteral converts a string ListMap key value into a fieldpath literal.
func stringLiteral(path fieldpath.Path, keyValue value.Value) (fieldpath.Literal, error) {
	if keyValue.Kind() != value.KindString {
		return fieldpath.Literal{}, keyKindMismatch(path, keyValue.Kind(), value.KindString)
	}

	stringPayload, _ := keyValue.String()

	return fieldpath.StringLiteral(stringPayload), nil
}
