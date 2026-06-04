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

package jsonconfig

import "testing"

func TestDefaultIsValid(t *testing.T) {
	if err := Validate(Default()); err != nil {
		t.Fatalf("Default() is invalid: %v", err)
	}
}

func TestResolveZeroConfigUsesDefaults(t *testing.T) {
	config, err := Resolve(Config{})
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}

	if config.Decode.Limits.MaxDepth != DefaultMaxDepth {
		t.Fatalf("decode max depth = %d; want %d", config.Decode.Limits.MaxDepth, DefaultMaxDepth)
	}
	if config.Decode.Limits.MaxNumberDigits != DefaultMaxNumberDigits {
		t.Fatalf("decode max number digits = %d; want %d", config.Decode.Limits.MaxNumberDigits, DefaultMaxNumberDigits)
	}
	if !config.Decode.Ownership.ValidateDocument {
		t.Fatalf("decode ownership validation = false; want true")
	}
	if config.Encode.Output.Layout != LayoutCompact {
		t.Fatalf("encode layout = %d; want compact", config.Encode.Output.Layout)
	}
	if config.Encode.Output.Indent != "  " {
		t.Fatalf("encode indent = %q; want two spaces", config.Encode.Output.Indent)
	}
	if config.Encode.Numbers.MaxDigits != DefaultEncodeMaxNumberDigits {
		t.Fatalf("encode max digits = %d; want %d", config.Encode.Numbers.MaxDigits, DefaultEncodeMaxNumberDigits)
	}
}

func TestResolveLeavesNoDefaultModes(t *testing.T) {
	config, err := Resolve(Config{})
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}

	assertNoDefaultModes(t, config)
}

func assertNoDefaultModes(t *testing.T, config Config) {
	t.Helper()

	if config.Decode.Numbers.Mode == NumberModeDefault {
		t.Fatalf("decode number mode still default")
	}
	if config.Decode.Strings.InvalidUTF8 == InvalidUTF8Default {
		t.Fatalf("decode invalid UTF-8 mode still default")
	}
	if config.Decode.Objects.DuplicateKeys == DuplicateKeyDefault {
		t.Fatalf("decode duplicate key mode still default")
	}
	if config.Decode.Objects.TrailingData == TrailingDataDefault {
		t.Fatalf("decode trailing data mode still default")
	}
	if config.Decode.Objects.UnknownEnvelopeFields == UnknownFieldDefault {
		t.Fatalf("decode unknown envelope field mode still default")
	}
	if config.Decode.Ownership.UnknownFields == UnknownFieldDefault {
		t.Fatalf("decode ownership unknown field mode still default")
	}
	if config.Encode.Output.Layout == LayoutDefault {
		t.Fatalf("encode layout mode still default")
	}
	if config.Encode.Output.FinalNewline == FinalNewlineDefault {
		t.Fatalf("encode final newline mode still default")
	}
	if config.Encode.Ordering.Mode == OrderingDefault {
		t.Fatalf("encode ordering mode still default")
	}
	if config.Encode.Strings.InvalidUTF8 == InvalidUTF8Default {
		t.Fatalf("encode invalid UTF-8 mode still default")
	}
	if config.Encode.Numbers.DecimalScale == DecimalScaleDefault {
		t.Fatalf("encode decimal scale mode still default")
	}
	if config.Encode.Numbers.FloatFormat == FloatFormatDefault {
		t.Fatalf("encode float format mode still default")
	}
	if config.Encode.Numbers.NegativeZero == NegativeZeroDefault {
		t.Fatalf("encode negative zero mode still default")
	}
	if config.Encode.Values.InvalidValue == InvalidValueDefault {
		t.Fatalf("encode invalid value mode still default")
	}
	if config.Encode.Values.Bytes == BytesEncodingDefault {
		t.Fatalf("encode bytes mode still default")
	}
	if config.Encode.Values.Timestamp == TimestampEncodingDefault {
		t.Fatalf("encode timestamp mode still default")
	}
	if config.Encode.Values.Date == DateEncodingDefault {
		t.Fatalf("encode date mode still default")
	}
	if config.Encode.Values.TimeOfDay == TimeOfDayEncodingDefault {
		t.Fatalf("encode time-of-day mode still default")
	}
	if config.Encode.Values.Duration == DurationEncodingDefault {
		t.Fatalf("encode duration mode still default")
	}
	if config.Encode.Object.TypeMeta == TypeMetaDefault {
		t.Fatalf("encode type meta mode still default")
	}
	if config.Encode.Object.Metadata == MetadataDefault {
		t.Fatalf("encode metadata mode still default")
	}
	if config.Encode.Object.Observed == ObservedDefault {
		t.Fatalf("encode observed mode still default")
	}
	if config.Encode.Ownership.Normalize == OwnershipNormalizeDefault {
		t.Fatalf("encode ownership normalize mode still default")
	}
	if config.Encode.Ownership.EmptyDesired == EmptyOwnershipSurfaceDefault {
		t.Fatalf("encode ownership empty desired mode still default")
	}
	if config.Encode.Ownership.EmptyEntries == EmptyEntriesDefault {
		t.Fatalf("encode ownership empty entries mode still default")
	}
}
