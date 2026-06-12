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

	"arcoris.dev/apimachinery/api/apidocument"
	"arcoris.dev/apimachinery/api/codec"
	"arcoris.dev/apimachinery/api/codecjson/jsonconfig"
)

// objectToNode converts a value-backed object envelope to stable JSON fields.
//
// Envelope member order is stable regardless of Deterministic. Deterministic
// controls only nested value objects and ownership state.
func objectToNode(path jsonPath, obj codec.Object, decode resolvedDecodeConfig, encode resolvedEncodeConfig) (jsonNode, error) {
	members := make([]jsonMember, 0, 5)
	var err error
	members, err = appendTypeMetaMembers(path, members, obj, encode)
	if err != nil {
		return jsonNode{}, err
	}

	if encode.metadata == jsonconfig.MetadataEmitEmpty || !obj.ObjectMeta.IsZero() {
		metadataNode, err := objectMetaToNode(path.Member(apidocument.ObjectFieldMetadata.String()), obj.ObjectMeta, decode)
		if err != nil {
			return jsonNode{}, err
		}
		members = append(members, jsonMember{name: apidocument.ObjectFieldMetadata.String(), value: metadataNode})
	}

	desiredNode, err := valueToNode(path.Member(apidocument.ObjectFieldDesired.String()), obj.Desired, encode)
	if err != nil {
		return jsonNode{}, err
	}
	members = append(members, jsonMember{name: apidocument.ObjectFieldDesired.String(), value: desiredNode})

	if obj.Observed != nil {
		observedNode, err := valueToNode(path.Member(apidocument.ObjectFieldObserved.String()), *obj.Observed, encode)
		if err != nil {
			return jsonNode{}, err
		}
		members = append(members, jsonMember{name: apidocument.ObjectFieldObserved.String(), value: observedNode})
	} else if encode.observed == jsonconfig.ObservedEmitNullWhenAbsent {
		members = append(members, jsonMember{
			name:  apidocument.ObjectFieldObserved.String(),
			value: jsonNode{kind: jsonKindNull},
		})
	}

	return jsonNode{kind: jsonKindObject, members: members}, nil
}

// appendTypeMetaMembers appends apiVersion and kind in canonical envelope order.
func appendTypeMetaMembers(
	path jsonPath,
	members []jsonMember,
	obj codec.Object,
	config resolvedEncodeConfig,
) ([]jsonMember, error) {
	if config.typeMeta == jsonconfig.TypeMetaRequire {
		switch {
		case obj.TypeMeta.APIVersion.IsZero():
			return nil, errorAt(path.Member(apidocument.ObjectFieldAPIVersion.String()), ErrInvalidEnvelope, errors.Join(codec.ErrEncodeFailed, codec.ErrInvalidDocument), ErrorReasonInvalidEnvelope, "apiVersion is required by JSON encode config")
		case obj.TypeMeta.Kind.IsZero():
			return nil, errorAt(path.Member(apidocument.ObjectFieldKind.String()), ErrInvalidEnvelope, errors.Join(codec.ErrEncodeFailed, codec.ErrInvalidDocument), ErrorReasonInvalidEnvelope, "kind is required by JSON encode config")
		}
	}
	if !obj.TypeMeta.APIVersion.IsZero() {
		members = append(members, jsonMember{
			name:  apidocument.ObjectFieldAPIVersion.String(),
			value: jsonNode{kind: jsonKindString, stringValue: obj.TypeMeta.APIVersion.String()},
		})
	}
	if !obj.TypeMeta.Kind.IsZero() {
		members = append(members, jsonMember{
			name:  apidocument.ObjectFieldKind.String(),
			value: jsonNode{kind: jsonKindString, stringValue: obj.TypeMeta.Kind.String()},
		})
	}

	return members, nil
}
