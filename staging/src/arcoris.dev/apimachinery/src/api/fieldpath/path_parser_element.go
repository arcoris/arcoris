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

// parseFieldElement decodes one dot-prefixed fixed-field step.
func (p *pathParser) parseFieldElement() (Element, error) {
	p.pos++

	if p.done() {
		return Element{}, p.syntaxError("field element is truncated")
	}

	if p.peek() == '"' {
		name, err := p.parseQuotedString()
		if err != nil {
			return Element{}, err
		}

		return FieldElementFromString(name)
	}

	start := p.pos

	if !p.consumeSimpleFieldStart() {
		return Element{}, p.syntaxError("field element name is invalid")
	}

	for !p.done() && isSimpleFieldContinue(p.peek()) {
		p.pos++
	}

	return FieldElementFromString(p.text[start:p.pos])
}

// parseBracketElement decodes one bracketed key, index, or selector step.
func (p *pathParser) parseBracketElement() (Element, error) {
	p.pos++

	if p.done() {
		return Element{}, p.syntaxError("bracket element is truncated")
	}

	var (
		element Element
		err     error
	)

	switch p.peek() {
	case '"':
		var key string

		key, err = p.parseQuotedString()
		if err == nil {
			element, err = KeyElementFromString(key)
		}
	case '{':
		var selector Selector

		selector, err = p.parseSelector()
		if err == nil {
			element, err = NewSelectorElement(selector)
		}
	default:
		element, err = p.parseIndexElement()
	}

	if err != nil {
		return Element{}, err
	}

	if !p.consumeByte(']') {
		return Element{}, p.syntaxError("bracket element is not closed")
	}

	return element, nil
}

// parseIndexElement decodes one non-negative decimal list index.
func (p *pathParser) parseIndexElement() (Element, error) {
	start := p.pos

	if p.done() || p.peek() < '0' || p.peek() > '9' {
		return Element{}, p.syntaxError("index element must start with a decimal digit")
	}

	for !p.done() && p.peek() >= '0' && p.peek() <= '9' {
		p.pos++
	}

	index, err := strconv.Atoi(p.text[start:p.pos])
	if err != nil {
		return Element{}, p.syntaxError("index element is out of range")
	}

	return NewIndexElement(index)
}
