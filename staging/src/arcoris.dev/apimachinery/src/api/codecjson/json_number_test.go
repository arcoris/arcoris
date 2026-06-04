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

package codecjson

import (
	"testing"

	"arcoris.dev/apimachinery/api/codec"
	"arcoris.dev/apimachinery/api/value"
)

func TestParseJSONNumberIntegerBoundaries(t *testing.T) {
	testCases := map[string]string{
		"0":                    "0",
		"-1":                   "-1",
		"9223372036854775807":  "9223372036854775807",
		"-9223372036854775808": "-9223372036854775808",
	}

	for input, want := range testCases {
		t.Run(input, func(t *testing.T) {
			got, err := parseJSONNumber(rootPath(), input, newTestCodec(t).decode)
			requireNoError(t, err)
			requireIntegerText(t, got, want)
		})
	}
}

func TestParseJSONNumberUnsignedBoundary(t *testing.T) {
	got, err := parseJSONNumber(rootPath(), "18446744073709551615", newTestCodec(t).decode)
	requireNoError(t, err)

	requireIntegerText(t, got, "18446744073709551615")
}

func TestParseJSONNumberRejectsBeyondUint64Integer(t *testing.T) {
	_, err := parseJSONNumber(rootPath(), "18446744073709551616", newTestCodec(t).decode)

	requireErrorIs(t, err, ErrInvalidNumber)
	requireErrorIs(t, err, codec.ErrInvalidNumber)
}

func TestParseJSONNumberDecimal(t *testing.T) {
	got, err := parseJSONNumber(rootPath(), "1.25", newTestCodec(t).decode)
	requireNoError(t, err)

	requireDecimalText(t, got, "1.25")
}

func TestParseJSONNumberExponentPositive(t *testing.T) {
	got, err := parseJSONNumber(rootPath(), "-1.20e2", newTestCodec(t).decode)
	requireNoError(t, err)

	requireKind(t, got, value.KindDecimal)
	requireDecimalText(t, got, "-120")
}

func TestParseJSONNumberExponentNegative(t *testing.T) {
	got, err := parseJSONNumber(rootPath(), "1e-3", newTestCodec(t).decode)
	requireNoError(t, err)

	requireDecimalText(t, got, "0.001")
}

func TestParseJSONNumberRejectsHugeExpansion(t *testing.T) {
	config := newTestCodec(t).decode
	config.maxNumberDigits = 4

	_, err := parseJSONNumber(rootPath(), "1e10", config)

	requireErrorIs(t, err, ErrInvalidNumber)
	requireErrorIs(t, err, codec.ErrInvalidNumber)
}

func TestParseJSONNumberDoesNotUseFloatPrecision(t *testing.T) {
	got, err := parseJSONNumber(rootPath(), "9007199254740993", newTestCodec(t).decode)
	requireNoError(t, err)

	requireIntegerText(t, got, "9007199254740993")
}
