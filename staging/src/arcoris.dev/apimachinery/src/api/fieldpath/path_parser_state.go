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

import "strings"

// done reports whether the parser has consumed all input.
func (p *pathParser) done() bool {
	return p.pos >= len(p.text)
}

// peek returns the current input byte.
func (p *pathParser) peek() byte {
	return p.text[p.pos]
}

// consumeByte advances over ch when it is the next byte.
func (p *pathParser) consumeByte(ch byte) bool {
	if p.done() || p.peek() != ch {
		return false
	}

	p.pos++
	return true
}

// consumeKeyword advances over keyword when it is the next token prefix.
func (p *pathParser) consumeKeyword(keyword string) bool {
	if !strings.HasPrefix(p.text[p.pos:], keyword) {
		return false
	}

	p.pos += len(keyword)
	return true
}

// consumeSimpleFieldStart advances over the first byte of a simple dot-form
// field name.
func (p *pathParser) consumeSimpleFieldStart() bool {
	if p.done() || !isSimpleFieldStart(p.peek()) {
		return false
	}

	p.pos++
	return true
}
