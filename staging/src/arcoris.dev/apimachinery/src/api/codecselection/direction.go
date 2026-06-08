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

// Direction identifies whether a binding is used for decode or encode.
type Direction uint8

const (
	// DirectionDecode selects codecs for reading encoded API documents.
	DirectionDecode Direction = iota + 1

	// DirectionEncode selects codecs for writing encoded API documents.
	DirectionEncode
)

// IsZero reports whether d is absent.
func (d Direction) IsZero() bool {
	return d == 0
}

// String returns stable diagnostic text for d.
func (d Direction) String() string {
	switch d {
	case DirectionDecode:
		return "decode"
	case DirectionEncode:
		return "encode"
	default:
		return ""
	}
}

// Validate checks that d is one of the supported selection directions.
func (d Direction) Validate() error {
	switch d {
	case DirectionDecode, DirectionEncode:
		return nil
	default:
		return errorfAt(
			pathDirection,
			ErrInvalidBinding,
			ErrorReasonInvalidBinding,
			"codec selection direction %d is not supported",
			d,
		)
	}
}
