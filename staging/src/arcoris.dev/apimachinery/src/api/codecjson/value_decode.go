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
	"arcoris.dev/apimachinery/api/codec"
	"arcoris.dev/apimachinery/api/value"
)

// nodeToValue converts a JSON node into the generic value algebra.
//
// This is the descriptor-agnostic mapping: JSON strings become value strings,
// JSON numbers are classified by source token text, and JSON object/array order
// is preserved. Descriptor-aware interpretations such as bytes or timestamps are left
// to future descriptor-aware codecs.
func nodeToValue(path jsonPath, node jsonNode, config resolvedDecodeConfig) (value.Value, error) {
	switch node.kind {
	case jsonKindNull:
		return value.NullValue(), nil
	case jsonKindBool:
		return value.BoolValue(node.boolValue), nil
	case jsonKindString:
		return value.StringValue(node.stringValue), nil
	case jsonKindNumber:
		return parseJSONNumber(path, node.numberText, config)
	case jsonKindArray:
		return nodeToListValue(path, node, config)
	case jsonKindObject:
		return nodeToObjectValue(path, node, config)
	default:
		return value.Value{}, errorAt(path, ErrInvalidJSON, codec.ErrDecodeFailed, ErrorReasonInvalidJSON, "invalid JSON node")
	}
}

// nodeToListValue converts an ordered JSON array into a value list.
//
// The list conversion walks with indexed JSON paths so nested diagnostics point
// at the exact array element that failed conversion.
func nodeToListValue(path jsonPath, node jsonNode, config resolvedDecodeConfig) (value.Value, error) {
	items := make([]value.Value, 0, len(node.items))
	for i, item := range node.items {
		converted, err := nodeToValue(path.Index(i), item, config)
		if err != nil {
			return value.Value{}, err
		}
		items = append(items, converted)
	}

	converted, err := value.ListValue(items...)
	if err != nil {
		return value.Value{}, wrapAt(path, ErrInvalidJSON, codec.ErrDecodeFailed, ErrorReasonInvalidJSON, "invalid JSON array value", err)
	}

	return converted, nil
}

// nodeToObjectValue converts an ordered JSON object into a value record.
//
// Duplicate names have already been rejected by the node parser, so this helper
// can preserve source member order without needing a lossy map intermediary.
func nodeToObjectValue(path jsonPath, node jsonNode, config resolvedDecodeConfig) (value.Value, error) {
	members := make([]value.RecordMember, 0, len(node.members))
	for _, member := range node.members {
		memberPath := path.Member(member.name)
		name, err := value.NewMemberName(member.name)
		if err != nil {
			return value.Value{}, wrapAt(memberPath, ErrInvalidJSON, codec.ErrDecodeFailed, ErrorReasonInvalidJSON, "invalid JSON object member name", err)
		}

		converted, err := nodeToValue(memberPath, member.value, config)
		if err != nil {
			return value.Value{}, err
		}
		members = append(members, value.NewRecordMember(name, converted))
	}

	converted, err := value.RecordValue(members...)
	if err != nil {
		return value.Value{}, wrapAt(path, ErrInvalidJSON, codec.ErrDecodeFailed, ErrorReasonInvalidJSON, "invalid JSON object value", err)
	}

	return converted, nil
}
