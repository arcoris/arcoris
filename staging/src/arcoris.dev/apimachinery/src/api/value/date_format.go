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

import "time"

// formatDate returns the diagnostic YYYY-MM-DD date form without using fmt.
//
// Date formatting is intentionally local to the value package and is not a wire
// codec. Keeping it byte-oriented also keeps the low-level value model free from
// fmt's reflection-heavy dependency tree.
func formatDate(year int, month time.Month, day int) string {
	out := make([]byte, 0, 10)
	out = appendPaddedSignedDecimal(out, year, 4)
	out = append(out, '-')
	out = appendPaddedSignedDecimal(out, int(month), 2)
	out = append(out, '-')
	out = appendPaddedSignedDecimal(out, day, 2)

	return string(out)
}

// formatDateMonth returns the diagnostic YYYY-MM form used in date errors.
func formatDateMonth(year int, month time.Month) string {
	out := make([]byte, 0, 7)
	out = appendPaddedSignedDecimal(out, year, 4)
	out = append(out, '-')
	out = appendPaddedSignedDecimal(out, int(month), 2)

	return string(out)
}
