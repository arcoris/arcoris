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

package objectquery

import (
	"encoding/base64"
	"math"
	"strconv"
	"strings"

	"arcoris.dev/apimachinery/api/value"
)

// canonicalValueKey produces a deterministic, collision-resistant sort key for
// literal sets. It is internal comparison machinery, not a wire encoding.
func canonicalValueKey(v value.Value) string {
	switch v.Kind() {
	case value.KindNull:
		return atom("null", "")
	case value.KindBool:
		payload, _ := v.AsBool()
		return atom("bool", strconv.FormatBool(payload))
	case value.KindString:
		payload, _ := v.AsString()
		return atom("string", payload)
	case value.KindBytes:
		payload, _ := v.AsBytes()
		return atom("bytes", base64.StdEncoding.EncodeToString(payload))
	case value.KindInteger:
		payload, _ := v.AsInteger()
		return atom("integer", payload.String())
	case value.KindFloat:
		payload, _ := v.AsFloat()
		if payload == 0 {
			payload = 0
		}
		return atom("float", strconv.FormatFloat(payload, 'g', -1, 64))
	case value.KindDecimal:
		payload, _ := v.AsDecimal()
		return atom("decimal", canonicalDecimal(payload))
	case value.KindTimestamp:
		payload, _ := v.AsTimestamp()
		return atom("timestamp", payload.UTC().Format("2006-01-02T15:04:05.999999999Z07:00"))
	case value.KindDate:
		payload, _ := v.AsDate()
		return atom("date", payload.String())
	case value.KindTimeOfDay:
		payload, _ := v.AsTimeOfDay()
		return atom("timeOfDay", payload.String())
	case value.KindDuration:
		payload, _ := v.AsDuration()
		return atom("duration", strconv.FormatInt(int64(payload), 10))
	case value.KindRecord:
		view, _ := v.AsRecord()
		parts := make([]string, 0, view.Len())
		view.ForEach(func(_ int, member value.RecordMember) bool {
			parts = append(parts, atom("member", member.Name.String())+canonicalValueKey(member.Value))
			return true
		})
		return atom("record", strings.Join(parts, ""))
	case value.KindList:
		view, _ := v.AsList()
		parts := make([]string, 0, view.Len())
		view.ForEach(func(_ int, nested value.Value) bool {
			parts = append(parts, canonicalValueKey(nested))
			return true
		})
		return atom("list", strings.Join(parts, ""))
	default:
		return atom("invalid", "")
	}
}

// atom length-prefixes payload so nested delimiters cannot collide.
func atom(kind string, payload string) string {
	return kind + "\x00" + payload + "\x00" + strconv.Itoa(len(payload))
}

// canonicalDecimal normalizes decimal representation to numeric equality.
func canonicalDecimal(decimal value.Decimal) string {
	if decimal.IsZero() {
		return "0"
	}

	coefficient := decimal.Coefficient()
	scale := decimal.Scale()
	for scale > 0 && strings.HasSuffix(coefficient, "0") {
		coefficient = strings.TrimSuffix(coefficient, "0")
		scale--
	}

	sign := "+"
	if decimal.IsNegative() {
		sign = "-"
	}
	if scale > math.MaxInt32 {
		return sign + coefficient + "e-" + strconv.FormatUint(uint64(scale), 10)
	}

	return sign + coefficient + "e-" + strconv.Itoa(int(scale))
}
