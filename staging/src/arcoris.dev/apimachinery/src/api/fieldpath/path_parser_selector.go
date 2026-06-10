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

// parseSelector decodes one associative-list selector body.
func (p *pathParser) parseSelector() (Selector, error) {
	if !p.consumeByte('{') {
		return Selector{}, p.syntaxError("selector must start with '{'")
	}

	entries := make([]SelectorEntry, 0, 2)

	for {
		if p.done() {
			return Selector{}, p.syntaxError("selector is truncated")
		}

		if p.peek() == '}' {
			p.pos++
			break
		}

		entry, err := p.parseSelectorEntry()
		if err != nil {
			return Selector{}, err
		}

		entries = append(entries, entry)

		if p.done() {
			return Selector{}, p.syntaxError("selector is truncated")
		}

		if p.peek() == '}' {
			p.pos++
			break
		}

		if !p.consumeByte(',') {
			return Selector{}, p.syntaxError("selector entries must be comma-separated")
		}
	}

	selector, err := NewSelector(entries...)
	if err != nil {
		return Selector{}, nested(
			ErrInvalidPath,
			ErrorReasonInvalidSelector,
			"selector text is invalid",
			err,
		)
	}

	return selector, nil
}

// parseSelectorEntry decodes one "field":literal pair inside a selector.
func (p *pathParser) parseSelectorEntry() (SelectorEntry, error) {
	field, err := p.parseQuotedString()
	if err != nil {
		return SelectorEntry{}, err
	}

	if !p.consumeByte(':') {
		return SelectorEntry{}, p.syntaxError("selector entry must contain ':'")
	}

	value, err := p.parseLiteral()
	if err != nil {
		return SelectorEntry{}, err
	}

	entry, err := SelectorEntryFromString(field, value)
	if err != nil {
		return SelectorEntry{}, nested(
			ErrInvalidPath,
			ErrorReasonInvalidEntry,
			"selector entry text is invalid",
			err,
		)
	}

	return entry, nil
}
