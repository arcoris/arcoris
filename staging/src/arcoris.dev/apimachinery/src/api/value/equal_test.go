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

package value

import "testing"

func TestEqualConcreteValues(t *testing.T) {
	left := MustRecordValue(
		MustRecordMember("name", StringValue("api")),
		MustRecordMember("items", MustListValue(StringValue("a,b"), StringValue("c=d"))),
	)
	right := MustRecordValue(
		MustRecordMember("name", StringValue("api")),
		MustRecordMember("items", MustListValue(StringValue("a,b"), StringValue("c=d"))),
	)

	if !Equal(left, right) {
		t.Fatal("Equal(record, record) = false; want true")
	}
	if Equal(left, MustRecordValue(MustRecordMember("name", StringValue("api")))) {
		t.Fatal("Equal(record, shorter record) = true; want false")
	}
}

func TestEqualDecimalUsesNumericSemantics(t *testing.T) {
	left := DecimalValue(MustParseDecimal("1.20"))
	right := DecimalValue(MustParseDecimal("1.2"))
	other := DecimalValue(MustParseDecimal("1.21"))

	if !Equal(left, right) {
		t.Fatal("Equal(1.20, 1.2) = false; want true")
	}
	if Equal(left, other) {
		t.Fatal("Equal(1.20, 1.21) = true; want false")
	}
}

func TestEqualFloatZeroSemantics(t *testing.T) {
	if !Equal(MustFloatValue(0), MustFloatValue(-0.0)) {
		t.Fatal("Equal(+0, -0) = false; want true")
	}
}
