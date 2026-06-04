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

package codecjson

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"slices"
	"strconv"
	"strings"

	"arcoris.dev/apimachinery/api/codec"
)

// encodeJSONDocument writes one JSON node to w.
func encodeJSONDocument(w io.Writer, node jsonNode, config nodeEncodeConfig) error {
	var buffer bytes.Buffer
	encoder := nodeEncoder{config: config}
	if err := encoder.writeNode(&buffer, node, 0); err != nil {
		return err
	}
	if config.finalNewline {
		buffer.WriteByte('\n')
	}
	if config.maxOutputBytes > 0 && int64(buffer.Len()) > config.maxOutputBytes {
		return errorfAt(
			rootPath(),
			ErrUnsupportedValue,
			codec.ErrEncodeFailed,
			ErrorReasonUnsupportedValue,
			"JSON output is %d bytes; maximum is %d",
			buffer.Len(),
			config.maxOutputBytes,
		)
	}
	_, err := w.Write(buffer.Bytes())
	return err
}

// nodeEncodeConfig stores JSON node writer behavior.
type nodeEncodeConfig struct {
	// pretty requests stable indentation.
	pretty bool

	// indent is the indentation unit for pretty output.
	indent string

	// finalNewline appends a newline after one document.
	finalNewline bool

	// escapeHTML mirrors encoding/json's HTML escaping policy for strings.
	escapeHTML bool

	// maxDepth is the inclusive JSON output nesting limit.
	maxDepth int

	// maxOutputBytes bounds encoded output bytes. Zero means unlimited.
	maxOutputBytes int64

	// maxNumberDigits bounds generated number text.
	maxNumberDigits int
}

// nodeEncoder writes jsonNode without map-based reordering.
type nodeEncoder struct {
	// config is immutable writer policy for this document.
	config nodeEncodeConfig
}

// writeNode writes one node at indentation depth.
func (e nodeEncoder) writeNode(out *bytes.Buffer, node jsonNode, depth int) error {
	if e.config.maxDepth > 0 && depth+1 > e.config.maxDepth {
		return errorAt(
			rootPath(),
			ErrUnsupportedValue,
			errors.Join(codec.ErrEncodeFailed, codec.ErrDepthExceeded),
			ErrorReasonUnsupportedValue,
			"JSON output exceeds maximum depth",
		)
	}

	switch node.kind {
	case jsonKindNull:
		out.WriteString("null")
	case jsonKindBool:
		out.WriteString(strconv.FormatBool(node.boolValue))
	case jsonKindString:
		out.WriteString(quoteJSONString(node.stringValue, e.config.escapeHTML))
	case jsonKindNumber:
		if e.config.maxNumberDigits > 0 && numberDigitCost(node.numberText) > e.config.maxNumberDigits {
			return errorAt(
				rootPath(),
				ErrInvalidNumber,
				errors.Join(codec.ErrEncodeFailed, codec.ErrInvalidNumber),
				ErrorReasonInvalidNumber,
				"JSON number exceeds maximum digit budget",
			)
		}
		out.WriteString(node.numberText)
	case jsonKindArray:
		return e.writeArray(out, node.items, depth)
	case jsonKindObject:
		return e.writeObject(out, node.members, depth)
	default:
		out.WriteString("null")
	}

	return nil
}

// writeArray writes array items in stored order.
func (e nodeEncoder) writeArray(out *bytes.Buffer, items []jsonNode, depth int) error {
	out.WriteByte('[')
	if len(items) == 0 {
		out.WriteByte(']')
		return nil
	}

	for i, item := range items {
		if i > 0 {
			out.WriteByte(',')
		}
		e.writePrettyLine(out, depth+1)
		if err := e.writeNode(out, item, depth+1); err != nil {
			return err
		}
	}
	e.writePrettyLine(out, depth)
	out.WriteByte(']')

	return nil
}

// writeObject writes object members in stored order.
func (e nodeEncoder) writeObject(out *bytes.Buffer, members []jsonMember, depth int) error {
	out.WriteByte('{')
	if len(members) == 0 {
		out.WriteByte('}')
		return nil
	}

	for i, member := range members {
		if i > 0 {
			out.WriteByte(',')
		}
		e.writePrettyLine(out, depth+1)
		out.WriteString(quoteJSONString(member.name, e.config.escapeHTML))
		if e.config.pretty {
			out.WriteString(": ")
		} else {
			out.WriteByte(':')
		}
		if err := e.writeNode(out, member.value, depth+1); err != nil {
			return err
		}
	}
	e.writePrettyLine(out, depth)
	out.WriteByte('}')

	return nil
}

// writePrettyLine writes newline and indentation when pretty output is enabled.
func (e nodeEncoder) writePrettyLine(out *bytes.Buffer, depth int) {
	if !e.config.pretty {
		return
	}

	out.WriteByte('\n')
	out.WriteString(strings.Repeat(e.config.indent, depth))
}

// quoteJSONString quotes s using the requested HTML escaping policy.
func quoteJSONString(s string, escapeHTML bool) string {
	if escapeHTML {
		data, _ := json.Marshal(s)
		return string(data)
	}

	var buffer bytes.Buffer
	encoder := json.NewEncoder(&buffer)
	encoder.SetEscapeHTML(false)
	_ = encoder.Encode(s)
	return strings.TrimSuffix(buffer.String(), "\n")
}

// sortedMembers returns members sorted by name for deterministic value objects.
func sortedMembers(members []jsonMember) []jsonMember {
	out := slices.Clone(members)
	slices.SortFunc(out, func(a jsonMember, b jsonMember) int {
		return strings.Compare(a.name, b.name)
	})

	return out
}
