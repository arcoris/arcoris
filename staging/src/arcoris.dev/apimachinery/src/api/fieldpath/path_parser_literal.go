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

import "strconv"

// parseLiteral decodes one supported selector literal.
//
// The first pass intentionally accepts only string, bool, and exact integer
// literals because selectors are stable identity keys, not general value trees.
func (p *pathParser) parseLiteral() (Literal, error) {
	if p.done() {
		return Literal{}, p.syntaxError("literal is truncated")
	}

	switch p.peek() {
	case '"':
		value, err := p.parseQuotedString()
		if err != nil {
			return Literal{}, err
		}

		return StringLiteral(value), nil
	case 't':
		if p.consumeKeyword("true") {
			return BoolLiteral(true), nil
		}
	case 'f':
		if p.consumeKeyword("false") {
			return BoolLiteral(false), nil
		}
	case '-':
		return p.parseSignedIntegerLiteral()
	default:
		if p.peek() >= '0' && p.peek() <= '9' {
			return p.parseUnsignedIntegerLiteral()
		}
	}

	return Literal{}, p.syntaxError("literal token is invalid")
}

// parseSignedIntegerLiteral decodes one exact int64 literal.
func (p *pathParser) parseSignedIntegerLiteral() (Literal, error) {
	start := p.pos
	p.pos++

	if p.done() || p.peek() < '0' || p.peek() > '9' {
		return Literal{}, p.syntaxError("signed integer literal is invalid")
	}

	for !p.done() && p.peek() >= '0' && p.peek() <= '9' {
		p.pos++
	}

	value, err := strconv.ParseInt(p.text[start:p.pos], 10, 64)
	if err != nil {
		return Literal{}, p.syntaxError("signed integer literal is out of range")
	}

	return Int64Literal(value), nil
}

// parseUnsignedIntegerLiteral decodes one exact uint64 literal.
func (p *pathParser) parseUnsignedIntegerLiteral() (Literal, error) {
	start := p.pos

	for !p.done() && p.peek() >= '0' && p.peek() <= '9' {
		p.pos++
	}

	value, err := strconv.ParseUint(p.text[start:p.pos], 10, 64)
	if err != nil {
		return Literal{}, p.syntaxError("unsigned integer literal is out of range")
	}

	return Uint64Literal(value), nil
}
