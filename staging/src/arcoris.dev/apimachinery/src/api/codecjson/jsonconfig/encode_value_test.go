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

func TestDefaultEncodeValueConfig(t *testing.T) {
	t.Parallel()

	config := defaultEncodeValueConfig()

	if config.InvalidValue != InvalidValueReject {
		t.Fatalf("invalid value = %d; want reject", config.InvalidValue)
	}
	if config.Bytes != BytesEncodingReject {
		t.Fatalf("bytes = %d; want reject", config.Bytes)
	}
	if config.Timestamp != TimestampEncodingReject {
		t.Fatalf("timestamp = %d; want reject", config.Timestamp)
	}
	if config.Date != DateEncodingReject {
		t.Fatalf("date = %d; want reject", config.Date)
	}
	if config.TimeOfDay != TimeOfDayEncodingReject {
		t.Fatalf("time of day = %d; want reject", config.TimeOfDay)
	}
	if config.Duration != DurationEncodingReject {
		t.Fatalf("duration = %d; want reject", config.Duration)
	}
}

func TestResolveEncodeValueConfig(t *testing.T) {
	t.Parallel()

	config := EncodeValueConfig{}
	resolveEncodeValueConfig(&config)

	if config.InvalidValue == InvalidValueDefault {
		t.Fatalf("invalid value still default")
	}
	if config.Bytes == BytesEncodingDefault {
		t.Fatalf("bytes still default")
	}
	if config.Timestamp == TimestampEncodingDefault {
		t.Fatalf("timestamp still default")
	}
	if config.Date == DateEncodingDefault {
		t.Fatalf("date still default")
	}
	if config.TimeOfDay == TimeOfDayEncodingDefault {
		t.Fatalf("time of day still default")
	}
	if config.Duration == DurationEncodingDefault {
		t.Fatalf("duration still default")
	}
}

func TestValidateEncodeValueConfigRejectsInvalidValues(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		config EncodeValueConfig
		target error
		path   string
	}{
		"invalid value unknown": {
			config: EncodeValueConfig{InvalidValue: InvalidValueMode(99), Bytes: BytesEncodingReject, Timestamp: TimestampEncodingReject, Date: DateEncodingReject, TimeOfDay: TimeOfDayEncodingReject, Duration: DurationEncodingReject},
			target: ErrInvalidConfig,
			path:   "encode.values.invalid_value",
		},
		"bytes unknown": {
			config: EncodeValueConfig{InvalidValue: InvalidValueReject, Bytes: BytesEncodingMode(99), Timestamp: TimestampEncodingReject, Date: DateEncodingReject, TimeOfDay: TimeOfDayEncodingReject, Duration: DurationEncodingReject},
			target: ErrInvalidConfig,
			path:   "encode.values.bytes",
		},
		"bytes unsupported": {
			config: EncodeValueConfig{InvalidValue: InvalidValueReject, Bytes: BytesEncodingBase64Std, Timestamp: TimestampEncodingReject, Date: DateEncodingReject, TimeOfDay: TimeOfDayEncodingReject, Duration: DurationEncodingReject},
			target: ErrUnsupportedConfig,
			path:   "encode.values.bytes",
		},
		"timestamp unsupported": {
			config: EncodeValueConfig{InvalidValue: InvalidValueReject, Bytes: BytesEncodingReject, Timestamp: TimestampEncodingRFC3339, Date: DateEncodingReject, TimeOfDay: TimeOfDayEncodingReject, Duration: DurationEncodingReject},
			target: ErrUnsupportedConfig,
			path:   "encode.values.timestamp",
		},
		"date unsupported": {
			config: EncodeValueConfig{InvalidValue: InvalidValueReject, Bytes: BytesEncodingReject, Timestamp: TimestampEncodingReject, Date: DateEncodingISO8601, TimeOfDay: TimeOfDayEncodingReject, Duration: DurationEncodingReject},
			target: ErrUnsupportedConfig,
			path:   "encode.values.date",
		},
		"time of day unsupported": {
			config: EncodeValueConfig{InvalidValue: InvalidValueReject, Bytes: BytesEncodingReject, Timestamp: TimestampEncodingReject, Date: DateEncodingReject, TimeOfDay: TimeOfDayEncodingISO8601, Duration: DurationEncodingReject},
			target: ErrUnsupportedConfig,
			path:   "encode.values.time_of_day",
		},
		"duration unsupported": {
			config: EncodeValueConfig{InvalidValue: InvalidValueReject, Bytes: BytesEncodingReject, Timestamp: TimestampEncodingReject, Date: DateEncodingReject, TimeOfDay: TimeOfDayEncodingReject, Duration: DurationEncodingISO8601},
			target: ErrUnsupportedConfig,
			path:   "encode.values.duration",
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			err := validateEncodeValueConfig(testCase.config)
			requireConfigErrorIs(t, err, testCase.target)
			requireErrorTextContains(t, err, testCase.path)
		})
	}
}
