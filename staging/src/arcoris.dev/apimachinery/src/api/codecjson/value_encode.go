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
	"errors"

	"arcoris.dev/apimachinery/api/codec"
	"arcoris.dev/apimachinery/api/value"
)

// valueToNode converts a concrete value into a JSON node.
//
// The conversion accepts only value kinds that can be decoded back by the
// generic JSON value codec without descriptors. Bytes and temporal values are
// rejected instead of being silently stringified into an ambiguous JSON string.
func valueToNode(path jsonPath, v value.Value, config resolvedEncodeConfig) (jsonNode, error) {
	switch v.Kind() {
	case value.KindNull:
		return jsonNode{kind: jsonKindNull}, nil
	case value.KindBool:
		payload, _ := v.AsBool()
		return jsonNode{kind: jsonKindBool, boolValue: payload}, nil
	case value.KindString:
		payload, _ := v.AsString()
		return jsonNode{kind: jsonKindString, stringValue: payload}, nil
	case value.KindInteger:
		payload, _ := v.AsInteger()
		return jsonNode{kind: jsonKindNumber, numberText: payload.String()}, nil
	case value.KindDecimal:
		payload, _ := v.AsDecimal()
		return jsonNode{kind: jsonKindNumber, numberText: payload.String()}, nil
	case value.KindFloat:
		payload, _ := v.AsFloat()
		text, err := finiteFloatText(path, payload, config)
		if err != nil {
			return jsonNode{}, err
		}
		return jsonNode{kind: jsonKindNumber, numberText: text}, nil
	case value.KindList:
		return listValueToNode(path, v, config)
	case value.KindRecord:
		return recordValueToNode(path, v, config)
	case value.KindInvalid,
		value.KindBytes,
		value.KindTimestamp,
		value.KindDate,
		value.KindTimeOfDay,
		value.KindDuration:
		return jsonNode{}, unsupportedValue(path, v.Kind())
	default:
		return jsonNode{}, unsupportedValue(path, v.Kind())
	}
}

// listValueToNode converts a value list into a JSON array.
//
// List order is semantic JSON order and is never changed by deterministic
// output. Deterministic only affects object member ordering.
func listValueToNode(path jsonPath, v value.Value, config resolvedEncodeConfig) (jsonNode, error) {
	list, _ := v.AsList()
	items := list.Items()
	nodes := make([]jsonNode, 0, len(items))
	for i, item := range items {
		node, err := valueToNode(path.Index(i), item, config)
		if err != nil {
			return jsonNode{}, err
		}
		nodes = append(nodes, node)
	}

	return jsonNode{kind: jsonKindArray, items: nodes}, nil
}

// recordValueToNode converts a value record into an ordered JSON object.
//
// Default output preserves value.RecordView member order. Deterministic output
// sorts the private JSON members after conversion so nested diagnostics still
// use the original semantic member path while output ordering becomes stable.
func recordValueToNode(path jsonPath, v value.Value, config resolvedEncodeConfig) (jsonNode, error) {
	record, _ := v.AsRecord()
	members := record.Members()
	jsonMembers := make([]jsonMember, 0, len(members))
	for _, member := range members {
		name := member.Name.String()
		node, err := valueToNode(path.Member(name), member.Value, config)
		if err != nil {
			return jsonNode{}, err
		}
		jsonMembers = append(jsonMembers, jsonMember{name: name, value: node})
	}
	if config.deterministic {
		jsonMembers = sortedMembers(jsonMembers)
	}

	return jsonNode{kind: jsonKindObject, members: jsonMembers}, nil
}

// unsupportedValue creates the explicit generic JSON unsupported-kind error.
func unsupportedValue(path jsonPath, kind value.Kind) error {
	return errorfAt(
		path,
		ErrUnsupportedValue,
		errors.Join(codec.ErrEncodeFailed, codec.ErrUnsupportedFeature),
		ErrorReasonUnsupportedValue,
		"value kind %s cannot round-trip through descriptor-agnostic JSON",
		kind,
	)
}
