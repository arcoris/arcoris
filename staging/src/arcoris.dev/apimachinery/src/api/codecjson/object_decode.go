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
	"arcoris.dev/apimachinery/api/apidocument"
	"arcoris.dev/apimachinery/api/codec"
	apiidentity "arcoris.dev/apimachinery/api/identity"
	"arcoris.dev/apimachinery/api/meta"
	"arcoris.dev/apimachinery/api/value"
)

// nodeToObject converts the JSON object envelope into the value-backed object.
//
// The conversion is intentionally resource-agnostic. It checks only the JSON
// envelope shape and lexical identity fields that are present in the document.
func nodeToObject(path jsonPath, node jsonNode, config resolvedDecodeConfig) (codec.Object, error) {
	fields := apidocument.Fields().Object()

	if err := requireObject(path, node, "object envelope root must be a JSON object"); err != nil {
		return codec.Object{}, err
	}
	if config.rejectUnknownEnvelopeFields {
		if err := rejectUnknownMembers(path, node, allowObjectEnvelopeField, "object envelope"); err != nil {
			return codec.Object{}, err
		}
	}

	typeMeta, err := nodeToTypeMeta(path, node)
	if err != nil {
		return codec.Object{}, err
	}
	objectMeta, err := nodeToOptionalObjectMeta(path, node, config)
	if err != nil {
		return codec.Object{}, err
	}
	desired, err := nodeToRequiredDesired(path, node, config)
	if err != nil {
		return codec.Object{}, err
	}

	out := codec.Object{
		TypeMeta:   typeMeta,
		ObjectMeta: objectMeta,
		Desired:    desired,
	}
	if observedNode, ok := node.member(fields.Observed().String()); ok {
		observed, err := nodeToValue(path.Member(fields.Observed().String()), observedNode, config)
		if err != nil {
			return codec.Object{}, err
		}
		out.Observed = &observed
	}

	return out, nil
}

// nodeToRequiredDesired extracts the desired payload while preserving explicit null.
func nodeToRequiredDesired(path jsonPath, node jsonNode, config resolvedDecodeConfig) (value.Value, error) {
	fields := apidocument.Fields().Object()
	desiredNode, ok := node.member(fields.Desired().String())
	if !ok {
		return value.Value{}, errorAt(
			path.Member(fields.Desired().String()),
			ErrInvalidEnvelope,
			codec.ErrInvalidDocument,
			ErrorReasonInvalidEnvelope,
			"object envelope desired field is required",
		)
	}

	return nodeToValue(path.Member(fields.Desired().String()), desiredNode, config)
}

// nodeToTypeMeta extracts optional apiVersion/kind fields without resource lookup.
func nodeToTypeMeta(path jsonPath, node jsonNode) (meta.TypeMeta, error) {
	fields := apidocument.Fields().Object()
	var typeMeta meta.TypeMeta
	if apiVersionNode, ok := node.member(fields.APIVersion().String()); ok {
		apiVersion, err := expectString(path.Member(fields.APIVersion().String()), apiVersionNode, "apiVersion must be a JSON string")
		if err != nil {
			return meta.TypeMeta{}, err
		}
		parsed, err := apiidentity.ParseGroupVersion(apiVersion)
		if err != nil {
			return meta.TypeMeta{}, wrapAt(
				path.Member(fields.APIVersion().String()),
				ErrInvalidEnvelope,
				codec.ErrInvalidDocument,
				ErrorReasonInvalidEnvelope,
				"apiVersion is not canonical",
				err,
			)
		}
		typeMeta.APIVersion = parsed
	}
	if kindNode, ok := node.member(fields.Kind().String()); ok {
		kind, err := expectString(path.Member(fields.Kind().String()), kindNode, "kind must be a JSON string")
		if err != nil {
			return meta.TypeMeta{}, err
		}
		parsed, err := apiidentity.ParseKind(kind)
		if err != nil {
			return meta.TypeMeta{}, wrapAt(
				path.Member(fields.Kind().String()),
				ErrInvalidEnvelope,
				codec.ErrInvalidDocument,
				ErrorReasonInvalidEnvelope,
				"kind is not canonical",
				err,
			)
		}
		typeMeta.Kind = parsed
	}

	return typeMeta, nil
}
