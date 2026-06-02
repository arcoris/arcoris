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

// signedIntegerLiteral converts a signed integer ListMap key value into a
// fieldpath literal.
func signedIntegerLiteral(
	path fieldpath.Path,
	keyValue value.Value,
) (fieldpath.Literal, error) {
	if keyValue.Kind() != value.KindInteger {
		return fieldpath.Literal{}, keyKindMismatch(path, keyValue.Kind(), value.KindInteger)
	}

	integerPayload, _ := keyValue.Integer()
	signedPayload, ok := integerPayload.Int64()
	if !ok {
		return fieldpath.Literal{}, failure(
			path,
			FailureKeyIntegerRange,
			"ListMap key integer does not fit signed selector type",
		)
	}

	return fieldpath.Int64Literal(signedPayload), nil
}

// unsignedIntegerLiteral converts an unsigned integer ListMap key value into a
// fieldpath literal.
func unsignedIntegerLiteral(
	path fieldpath.Path,
	keyValue value.Value,
) (fieldpath.Literal, error) {
	if keyValue.Kind() != value.KindInteger {
		return fieldpath.Literal{}, keyKindMismatch(path, keyValue.Kind(), value.KindInteger)
	}

	integerPayload, _ := keyValue.Integer()
	unsignedPayload, ok := integerPayload.Uint64()
	if !ok {
		return fieldpath.Literal{}, failure(
			path,
			FailureKeyIntegerRange,
			"ListMap key integer does not fit unsigned selector type",
		)
	}

	return fieldpath.Uint64Literal(unsignedPayload), nil
}
