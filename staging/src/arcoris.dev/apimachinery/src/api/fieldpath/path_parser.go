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

// pathParser incrementally decodes the canonical Path.String grammar.
//
// The parser intentionally works at the byte level because the field-path text
// grammar is ASCII-structured. Quoted names and string literals are delegated
// to strconv.Unquote for escaping and UTF-8 validation in narrower helpers.
type pathParser struct {
	text string
	pos  int
}

// newPathParser prepares one parser instance for text.
func newPathParser(text string) pathParser {
	return pathParser{text: text}
}
