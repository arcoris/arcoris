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

// EncodeValueConfig controls descriptor-agnostic value encoding.
type EncodeValueConfig struct {
	// InvalidValue controls value.KindInvalid handling.
	InvalidValue InvalidValueMode

	// Bytes controls value.KindBytes handling.
	Bytes BytesEncodingMode

	// Timestamp controls value.KindTimestamp handling.
	Timestamp TimestampEncodingMode

	// Date controls value.KindDate handling.
	Date DateEncodingMode

	// TimeOfDay controls value.KindTimeOfDay handling.
	TimeOfDay TimeOfDayEncodingMode

	// Duration controls value.KindDuration handling.
	Duration DurationEncodingMode
}

// defaultEncodeValueConfig returns safe descriptor-agnostic value policy.
func defaultEncodeValueConfig() EncodeValueConfig {
	return EncodeValueConfig{
		InvalidValue: InvalidValueReject,
		Bytes:        BytesEncodingReject,
		Timestamp:    TimestampEncodingReject,
		Date:         DateEncodingReject,
		TimeOfDay:    TimeOfDayEncodingReject,
		Duration:     DurationEncodingReject,
	}
}

// resolveEncodeValueConfig applies value encode defaults in place.
func resolveEncodeValueConfig(config *EncodeValueConfig) {
	if config.InvalidValue == InvalidValueDefault {
		config.InvalidValue = InvalidValueReject
	}
	if config.Bytes == BytesEncodingDefault {
		config.Bytes = BytesEncodingReject
	}
	if config.Timestamp == TimestampEncodingDefault {
		config.Timestamp = TimestampEncodingReject
	}
	if config.Date == DateEncodingDefault {
		config.Date = DateEncodingReject
	}
	if config.TimeOfDay == TimeOfDayEncodingDefault {
		config.TimeOfDay = TimeOfDayEncodingReject
	}
	if config.Duration == DurationEncodingDefault {
		config.Duration = DurationEncodingReject
	}
}

// validateEncodeValueConfig checks generic value kind output policy.
func validateEncodeValueConfig(config EncodeValueConfig) error {
	switch {
	case !isKnownInvalidValueMode(config.InvalidValue):
		return invalidConfig("encode.values.invalid_value", "unknown invalid value mode %d", config.InvalidValue)
	case config.InvalidValue != InvalidValueReject:
		return unsupportedConfig("encode.values.invalid_value", "invalid value mode %d is not implemented", config.InvalidValue)
	case !isKnownBytesEncodingMode(config.Bytes):
		return invalidConfig("encode.values.bytes", "unknown bytes encoding mode %d", config.Bytes)
	case config.Bytes != BytesEncodingReject:
		return unsupportedConfig("encode.values.bytes", "bytes encoding mode %d is not implemented by generic codecjson", config.Bytes)
	case !isKnownTimestampEncodingMode(config.Timestamp):
		return invalidConfig("encode.values.timestamp", "unknown timestamp encoding mode %d", config.Timestamp)
	case config.Timestamp != TimestampEncodingReject:
		return unsupportedConfig("encode.values.timestamp", "timestamp encoding mode %d is not implemented by generic codecjson", config.Timestamp)
	case !isKnownDateEncodingMode(config.Date):
		return invalidConfig("encode.values.date", "unknown date encoding mode %d", config.Date)
	case config.Date != DateEncodingReject:
		return unsupportedConfig("encode.values.date", "date encoding mode %d is not implemented by generic codecjson", config.Date)
	case !isKnownTimeOfDayEncodingMode(config.TimeOfDay):
		return invalidConfig("encode.values.time_of_day", "unknown time-of-day encoding mode %d", config.TimeOfDay)
	case config.TimeOfDay != TimeOfDayEncodingReject:
		return unsupportedConfig("encode.values.time_of_day", "time-of-day encoding mode %d is not implemented by generic codecjson", config.TimeOfDay)
	case !isKnownDurationEncodingMode(config.Duration):
		return invalidConfig("encode.values.duration", "unknown duration encoding mode %d", config.Duration)
	case config.Duration != DurationEncodingReject:
		return unsupportedConfig("encode.values.duration", "duration encoding mode %d is not implemented by generic codecjson", config.Duration)
	default:
		return nil
	}
}
