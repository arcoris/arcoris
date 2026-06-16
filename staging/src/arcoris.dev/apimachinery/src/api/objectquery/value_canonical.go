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
	"strconv"
	"strings"

	"arcoris.dev/apimachinery/api/value"
)

// canonicalValueKey produces a deterministic equality key for value.Value.
// It is internal comparison machinery, not a wire or storage encoding.
func canonicalValueKey(v value.Value) string {
	switch v.Kind() {
	case value.KindNull:
		return "null:"
	case value.KindBool:
		payload, _ := v.AsBool()
		return "bool:" + strconv.FormatBool(payload)
	case value.KindString:
		payload, _ := v.AsString()
		return "string:" + payload
	case value.KindBytes:
		payload, _ := v.AsBytes()
		return "bytes:" + base64.StdEncoding.EncodeToString(payload)
	case value.KindInteger:
		payload, _ := v.AsInteger()
		return "integer:" + payload.String()
	case value.KindFloat:
		payload, _ := v.AsFloat()
		return "float:" + strconv.FormatFloat(payload, 'g', -1, 64)
	case value.KindDecimal:
		payload, _ := v.AsDecimal()
		return "decimal:" + payload.String()
	case value.KindTimestamp:
		payload, _ := v.AsTimestamp()
		return "timestamp:" + payload.UTC().Format("2006-01-02T15:04:05.999999999Z07:00")
	case value.KindDate:
		payload, _ := v.AsDate()
		return "date:" + payload.String()
	case value.KindTimeOfDay:
		payload, _ := v.AsTimeOfDay()
		return "timeOfDay:" + payload.String()
	case value.KindDuration:
		payload, _ := v.AsDuration()
		return "duration:" + payload.String()
	case value.KindRecord:
		view, _ := v.AsRecord()
		parts := make([]string, 0, view.Len())
		view.ForEach(func(_ int, member value.RecordMember) bool {
			parts = append(parts, member.Name.String()+"="+canonicalValueKey(member.Value))
			return true
		})
		return "record:" + strings.Join(parts, ",")
	case value.KindList:
		view, _ := v.AsList()
		parts := make([]string, 0, view.Len())
		view.ForEach(func(_ int, nested value.Value) bool {
			parts = append(parts, canonicalValueKey(nested))
			return true
		})
		return "list:" + strings.Join(parts, ",")
	default:
		return "invalid:"
	}
}
