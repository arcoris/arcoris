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

// parsePath reads one complete semantic path from the configured text.
func (p *pathParser) parsePath() (Path, error) {
	if !p.consumeByte('$') {
		return Path{}, p.syntaxError("path must start with '$'")
	}

	elements := make([]Element, 0, 4)

	for !p.done() {
		switch p.peek() {
		case '.':
			element, err := p.parseFieldElement()
			if err != nil {
				return Path{}, err
			}

			elements = append(elements, element)
		case '[':
			element, err := p.parseBracketElement()
			if err != nil {
				return Path{}, err
			}

			elements = append(elements, element)
		default:
			return Path{}, p.syntaxError("unexpected token in path")
		}
	}

	path := Path{elements: elements}
	if err := path.Validate(); err != nil {
		return Path{}, err
	}

	return path, nil
}
