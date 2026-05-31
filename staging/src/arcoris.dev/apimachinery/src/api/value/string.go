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

// StringValue constructs a KindString Value.
//
// It stores Go string data as supplied. UTF-8, length, and semantic-token rules
// belong to codecs or future descriptor-aware validation. This keeps the value
// model focused on storing payload data rather than enforcing schema rules.
func StringValue(v string) Value {
	return Value{kind: KindString, stringValue: v}
}
