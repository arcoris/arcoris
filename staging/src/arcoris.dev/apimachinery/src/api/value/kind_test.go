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

func TestKindClassification(t *testing.T) {
	tests := []struct {
		kind      Kind
		text      string
		valid     bool
		scalar    bool
		number    bool
		temporal  bool
		composite bool
	}{
		{kind: KindInvalid, text: "invalid"},
		{kind: KindNull, text: "null", valid: true, scalar: true},
		{kind: KindBool, text: "bool", valid: true, scalar: true},
		{kind: KindString, text: "string", valid: true, scalar: true},
		{kind: KindBytes, text: "bytes", valid: true, scalar: true},
		{kind: KindInteger, text: "integer", valid: true, scalar: true, number: true},
		{kind: KindFloat, text: "float", valid: true, scalar: true, number: true},
		{kind: KindDecimal, text: "decimal", valid: true, scalar: true, number: true},
		{kind: KindTimestamp, text: "timestamp", valid: true, scalar: true, temporal: true},
		{kind: KindDate, text: "date", valid: true, scalar: true, temporal: true},
		{kind: KindTimeOfDay, text: "timeOfDay", valid: true, scalar: true, temporal: true},
		{kind: KindDuration, text: "duration", valid: true, scalar: true, temporal: true},
		{kind: KindRecord, text: "record", valid: true, composite: true},
		{kind: KindList, text: "list", valid: true, composite: true},
	}

	for _, tt := range tests {
		t.Run(tt.text, func(t *testing.T) {
			requireEqual(t, tt.kind.String(), tt.text)
			requireEqual(t, tt.kind.IsValid(), tt.valid)
			requireEqual(t, tt.kind.IsScalar(), tt.scalar)
			requireEqual(t, tt.kind.IsNumber(), tt.number)
			requireEqual(t, tt.kind.IsTemporal(), tt.temporal)
			requireEqual(t, tt.kind.IsComposite(), tt.composite)
		})
	}
}

func TestKindUnknownString(t *testing.T) {
	requireEqual(t, Kind(255).String(), "unknown")
	requireEqual(t, Kind(255).IsValid(), false)
}

func TestKindRecordIsOnlyKeyedPayloadKind(t *testing.T) {
	requireEqual(t, KindRecord.String(), "record")
	requireEqual(t, KindRecord.IsComposite(), true)
	requireEqual(t, KindList.IsComposite(), true)
	requireEqual(t, Kind(KindList+1).IsValid(), false)
}
