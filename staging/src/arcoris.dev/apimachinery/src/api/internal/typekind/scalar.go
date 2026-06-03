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

package typekind

import (
	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
)

// Scalar maps non-null scalar descriptor codes to their concrete payload kind.
//
// Integer width and signedness are descriptor constraints; concrete payloads use
// one integer kind for both signed and unsigned integer descriptors. TypeNull is
// intentionally not included because null is a presence-sensitive leaf in the
// caller packages, not part of this reusable scalar-kind table.
func Scalar(code types.TypeCode) (value.Kind, bool) {
	switch code {
	case types.TypeBool:
		return value.KindBool, true
	case types.TypeString:
		return value.KindString, true
	case types.TypeBytes:
		return value.KindBytes, true
	case types.TypeInt8,
		types.TypeInt16,
		types.TypeInt32,
		types.TypeInt64,
		types.TypeUint8,
		types.TypeUint16,
		types.TypeUint32,
		types.TypeUint64:
		return value.KindInteger, true
	case types.TypeFloat32,
		types.TypeFloat64:
		return value.KindFloat, true
	case types.TypeDecimal:
		return value.KindDecimal, true
	case types.TypeTimestamp:
		return value.KindTimestamp, true
	case types.TypeDate:
		return value.KindDate, true
	case types.TypeTime:
		return value.KindTimeOfDay, true
	case types.TypeDuration:
		return value.KindDuration, true
	default:
		return value.KindInvalid, false
	}
}
