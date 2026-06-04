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
	"unicode/utf8"

	"arcoris.dev/apimachinery/api/codec"
)

// decodeJSONDocument decodes exactly one JSON document into the ordered node model.
func decodeJSONDocument(r io.Reader, config resolvedDecodeConfig) (jsonNode, error) {
	data, err := readJSONInput(r, config)
	if err != nil {
		return jsonNode{}, err
	}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()

	decode := nodeDecoder{
		decoder:        decoder,
		maxDepth:       config.maxDepth,
		maxStringBytes: config.maxStringBytes,
	}

	node, err := decode.decodeValue(rootPath(), 1)
	if err != nil {
		return jsonNode{}, err
	}
	if err := rejectTrailingData(decoder); err != nil {
		return jsonNode{}, err
	}

	return node, nil
}

// readJSONInput reads one document and rejects raw invalid UTF-8 before parsing.
//
// JSON is a Unicode text format. encoding/json accepts invalid UTF-8 inside
// strings by replacing bytes with RuneError, which would lose source fidelity,
// so codecjson rejects such input at the byte boundary.
func readJSONInput(r io.Reader, config resolvedDecodeConfig) ([]byte, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, wrapAt(rootPath(), ErrInvalidJSON, codec.ErrDecodeFailed, ErrorReasonInvalidJSON, "JSON input cannot be read", err)
	}
	if config.maxDocumentBytes > 0 && int64(len(data)) > config.maxDocumentBytes {
		return nil, errorfAt(
			rootPath(),
			ErrInvalidJSON,
			codec.ErrDecodeFailed,
			ErrorReasonInvalidJSON,
			"JSON document is %d bytes; maximum is %d",
			len(data),
			config.maxDocumentBytes,
		)
	}
	if !utf8.Valid(data) {
		return nil, errorAt(rootPath(), ErrInvalidJSON, codec.ErrDecodeFailed, ErrorReasonInvalidJSON, "JSON input is not valid UTF-8")
	}

	return data, nil
}

// nodeDecoder owns token-based recursive JSON decoding state.
type nodeDecoder struct {
	// decoder is the token-mode JSON decoder for one input document.
	decoder *json.Decoder

	// maxDepth is the inclusive maximum document nesting depth.
	maxDepth int

	// maxStringBytes bounds decoded string token length. Zero means unlimited.
	maxStringBytes int
}

// decodeValue decodes one JSON token at path and depth.
func (d nodeDecoder) decodeValue(path jsonPath, depth int) (jsonNode, error) {
	if depth > d.maxDepth {
		return jsonNode{}, errorAt(
			path,
			ErrInvalidJSON,
			codec.ErrDepthExceeded,
			ErrorReasonMaxDepthExceeded,
			"JSON document exceeds maximum depth",
		)
	}

	token, err := d.decoder.Token()
	if err != nil {
		return jsonNode{}, wrapAt(path, ErrInvalidJSON, codec.ErrDecodeFailed, ErrorReasonInvalidJSON, "invalid JSON", err)
	}

	switch typed := token.(type) {
	case nil:
		return jsonNode{kind: jsonKindNull}, nil
	case bool:
		return jsonNode{kind: jsonKindBool, boolValue: typed}, nil
	case string:
		if err := d.rejectOversizedString(path, typed); err != nil {
			return jsonNode{}, err
		}
		return jsonNode{kind: jsonKindString, stringValue: typed}, nil
	case json.Number:
		return jsonNode{kind: jsonKindNumber, numberText: typed.String()}, nil
	case json.Delim:
		switch typed {
		case '[':
			return d.decodeArray(path, depth)
		case '{':
			return d.decodeObject(path, depth)
		default:
			return jsonNode{}, errorfAt(path, ErrInvalidJSON, codec.ErrDecodeFailed, ErrorReasonInvalidJSON, "unexpected JSON delimiter %q", typed)
		}
	default:
		return jsonNode{}, errorfAt(path, ErrInvalidJSON, codec.ErrDecodeFailed, ErrorReasonInvalidJSON, "unexpected JSON token %T", token)
	}
}

// decodeArray decodes an array after the opening delimiter has been consumed.
func (d nodeDecoder) decodeArray(path jsonPath, depth int) (jsonNode, error) {
	items := []jsonNode{}
	for index := 0; d.decoder.More(); index++ {
		item, err := d.decodeValue(path.Index(index), depth+1)
		if err != nil {
			return jsonNode{}, err
		}
		items = append(items, item)
	}
	if err := d.expectDelim(path, ']'); err != nil {
		return jsonNode{}, err
	}

	return jsonNode{kind: jsonKindArray, items: items}, nil
}

// decodeObject decodes an object after the opening delimiter has been consumed.
func (d nodeDecoder) decodeObject(path jsonPath, depth int) (jsonNode, error) {
	members := []jsonMember{}
	seen := map[string]struct{}{}
	for d.decoder.More() {
		keyToken, err := d.decoder.Token()
		if err != nil {
			return jsonNode{}, wrapAt(path, ErrInvalidJSON, codec.ErrDecodeFailed, ErrorReasonInvalidJSON, "invalid JSON object key", err)
		}
		key, ok := keyToken.(string)
		if !ok {
			return jsonNode{}, errorfAt(path, ErrInvalidJSON, codec.ErrDecodeFailed, ErrorReasonInvalidJSON, "expected JSON object key, got %T", keyToken)
		}
		if err := d.rejectOversizedString(path, key); err != nil {
			return jsonNode{}, err
		}
		if _, ok := seen[key]; ok {
			return jsonNode{}, errorfAt(path, ErrDuplicateKey, codec.ErrInvalidDocument, ErrorReasonDuplicateKey, "duplicate JSON object key %q", key)
		}
		seen[key] = struct{}{}

		value, err := d.decodeValue(path.Member(key), depth+1)
		if err != nil {
			return jsonNode{}, err
		}
		members = append(members, jsonMember{name: key, value: value})
	}
	if err := d.expectDelim(path, '}'); err != nil {
		return jsonNode{}, err
	}

	return jsonNode{kind: jsonKindObject, members: members}, nil
}

// rejectOversizedString enforces decoded string byte limits for keys and values.
func (d nodeDecoder) rejectOversizedString(path jsonPath, text string) error {
	if d.maxStringBytes == 0 || len(text) <= d.maxStringBytes {
		return nil
	}

	return errorfAt(
		path,
		ErrInvalidJSON,
		codec.ErrDecodeFailed,
		ErrorReasonInvalidJSON,
		"JSON string is %d bytes; maximum is %d",
		len(text),
		d.maxStringBytes,
	)
}

// expectDelim consumes the required closing delimiter.
func (d nodeDecoder) expectDelim(path jsonPath, want json.Delim) error {
	token, err := d.decoder.Token()
	if err != nil {
		return wrapAt(path, ErrInvalidJSON, codec.ErrDecodeFailed, ErrorReasonInvalidJSON, "invalid JSON delimiter", err)
	}
	if got, ok := token.(json.Delim); !ok || got != want {
		return errorfAt(path, ErrInvalidJSON, codec.ErrDecodeFailed, ErrorReasonInvalidJSON, "expected JSON delimiter %q", want)
	}

	return nil
}

// rejectTrailingData rejects any token after the first document.
func rejectTrailingData(decoder *json.Decoder) error {
	token, err := decoder.Token()
	if errors.Is(err, io.EOF) {
		return nil
	}
	if err != nil {
		return wrapAt(rootPath(), ErrInvalidJSON, codec.ErrDecodeFailed, ErrorReasonInvalidJSON, "invalid trailing JSON data", err)
	}

	return errorfAt(rootPath(), ErrTrailingData, codec.ErrDecodeFailed, ErrorReasonTrailingData, "unexpected trailing JSON token %v", token)
}
