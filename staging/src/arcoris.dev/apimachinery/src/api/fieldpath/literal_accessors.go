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

// BoolValue returns the stored boolean when l is a bool literal.
func (l Literal) BoolValue() (bool, bool) {
	if l.kind != LiteralBool {
		return false, false
	}

	return l.boolValue, true
}

// StringValue returns the stored string when l is a string literal.
func (l Literal) StringValue() (string, bool) {
	if l.kind != LiteralString {
		return "", false
	}

	return l.stringValue, true
}

// Int64Value returns the stored integer when l fits into int64.
func (l Literal) Int64Value() (int64, bool) {
	if l.kind != LiteralInteger {
		return 0, false
	}

	return l.intValue.int64Value()
}

// Uint64Value returns the stored integer when l is non-negative.
func (l Literal) Uint64Value() (uint64, bool) {
	if l.kind != LiteralInteger {
		return 0, false
	}

	return l.intValue.uint64Value()
}
