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

// ParseCanonical reconstructs one semantic path from the canonical text grammar.
//
// The parser is intentionally strict. It accepts the grammar emitted by
// Path.CanonicalText and rejects query syntax, wildcards, and non-canonical
// forms that would blur the path model's structural semantics.
//
// ParseCanonical is meant for diagnostics, tests, and future persistence helpers
// that need a lossless round-trip for the package's canonical text form. It is
// not a general query language and does not accept wildcard or filter syntax.
func ParseCanonical(text string) (Path, error) {
	p := newPathParser(text)

	path, err := p.parsePath()
	if err != nil {
		return Path{}, err
	}

	if canonical := path.CanonicalText(); canonical != text {
		return Path{}, nested(
			ErrInvalidPath,
			ErrorReasonNonCanonicalText,
			"path text is not canonical",
			ErrNonCanonicalText,
		)
	}

	return path, nil
}
