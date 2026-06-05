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

package lexical

// IsASCIILower reports whether b is an ASCII lowercase letter.
func IsASCIILower(b byte) bool { return b >= 'a' && b <= 'z' }

// IsASCIIUpper reports whether b is an ASCII uppercase letter.
func IsASCIIUpper(b byte) bool { return b >= 'A' && b <= 'Z' }

// IsASCIIDigit reports whether b is an ASCII decimal digit.
func IsASCIIDigit(b byte) bool { return b >= '0' && b <= '9' }

// IsASCIIAlpha reports whether b is an ASCII letter.
func IsASCIIAlpha(b byte) bool { return IsASCIILower(b) || IsASCIIUpper(b) }

// IsASCIIAlnum reports whether b is an ASCII letter or digit.
func IsASCIIAlnum(b byte) bool { return IsASCIIAlpha(b) || IsASCIIDigit(b) }

// IsDNS1123LabelChar reports whether b can appear inside a DNS-1123 label.
func IsDNS1123LabelChar(b byte) bool {
	return IsASCIILower(b) || IsASCIIDigit(b) || b == '-'
}

// IsDNS1123LabelEdge reports whether b can start or end a DNS-1123 label.
func IsDNS1123LabelEdge(b byte) bool {
	return IsASCIILower(b) || IsASCIIDigit(b)
}
