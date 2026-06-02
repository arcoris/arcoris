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
	"fmt"

	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
)

// literalFromValue converts one declared ListMap key value into a selector
// literal.
func literalFromValue(
	path fieldpath.Path,
	keyValue value.Value,
	descriptor types.Type,
	references referenceResolver,
	depth int,
) (fieldpath.Literal, error) {
	resolvedDescriptor, err := references.resolve(path, descriptor, depth)
	if err != nil {
		return fieldpath.Literal{}, err
	}

	switch resolvedDescriptor.Code() {
	case types.TypeBool:
		return boolLiteral(path, keyValue)
	case types.TypeString:
		return stringLiteral(path, keyValue)
	case types.TypeInt8,
		types.TypeInt16,
		types.TypeInt32,
		types.TypeInt64:
		return signedIntegerLiteral(path, keyValue)
	case types.TypeUint8,
		types.TypeUint16,
		types.TypeUint32,
		types.TypeUint64:
		return unsignedIntegerLiteral(path, keyValue)
	default:
		return fieldpath.Literal{}, failure(
			path,
			FailureInvalidDescriptor,
			fmt.Sprintf(
				"descriptor %s cannot provide ListMap key identity",
				resolvedDescriptor.Code(),
			),
		)
	}
}
