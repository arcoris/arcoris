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
	descriptor types.Descriptor,
	references referenceResolver,
	depth int,
) (fieldpath.Literal, error) {
	resolvedDescriptor, err := references.resolve(path, descriptor, depth)
	if err != nil {
		return fieldpath.Literal{}, err
	}

	switch resolvedDescriptor.Code() {
	case types.DescriptorBool:
		return boolLiteral(path, keyValue)
	case types.DescriptorString:
		return stringLiteral(path, keyValue)
	case types.DescriptorInt8,
		types.DescriptorInt16,
		types.DescriptorInt32,
		types.DescriptorInt64:
		return signedIntegerLiteral(path, keyValue)
	case types.DescriptorUint8,
		types.DescriptorUint16,
		types.DescriptorUint32,
		types.DescriptorUint64:
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
