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

package fieldpath

import (
	"strconv"
	"unicode/utf8"
)

// parseQuotedString decodes one JSON-style quoted string token.
func (p *pathParser) parseQuotedString() (string, error) {
	if !p.consumeByte('"') {
		return "", p.syntaxError("quoted string must start with '\"'")
	}

	start := p.pos - 1
	escaped := false

	for !p.done() {
		ch := p.peek()
		p.pos++

		if escaped {
			escaped = false
			continue
		}

		if ch == '\\' {
			escaped = true
			continue
		}

		if ch == '"' {
			quoted := p.text[start:p.pos]

			value, err := strconv.Unquote(quoted)
			if err != nil {
				return "", p.syntaxError("quoted string escape sequence is invalid")
			}

			if !utf8.ValidString(value) {
				return "", p.syntaxError("quoted string contains invalid UTF-8")
			}

			return value, nil
		}
	}

	return "", p.syntaxError("quoted string is not closed")
}
