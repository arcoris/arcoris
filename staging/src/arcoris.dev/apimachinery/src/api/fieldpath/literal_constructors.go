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

package fieldpath

// BoolLiteral constructs a boolean selector literal.
//
// The resulting literal is always valid because bool selectors carry no extra
// grammar beyond their kind.
func BoolLiteral(value bool) Literal {
	return Literal{
		kind:      LiteralBool,
		boolValue: value,
	}
}

// StringLiteral constructs a string selector literal.
//
// Empty strings are allowed at the fieldpath layer. Higher descriptor-aware
// validation may still reject them for specific associative-list keys.
func StringLiteral(value string) Literal {
	return Literal{
		kind:        LiteralString,
		stringValue: value,
	}
}

// Int64Literal constructs an exact signed integer selector literal.
//
// The private sign/magnitude representation preserves the full int64 domain,
// including math.MinInt64, without overflow tricks leaking into callers.
func Int64Literal(value int64) Literal {
	return Literal{
		kind:     LiteralInteger,
		intValue: newInt64(value),
	}
}

// Uint64Literal constructs an exact unsigned integer selector literal.
//
// Unsigned selector literals remain comparable against signed literals through
// one shared exact integer representation.
func Uint64Literal(value uint64) Literal {
	return Literal{
		kind:     LiteralInteger,
		intValue: newUint64(value),
	}
}
