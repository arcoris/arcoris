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

	"arcoris.dev/apimachinery/api/apidocument"
	"arcoris.dev/apimachinery/api/codec"
	"arcoris.dev/apimachinery/api/meta"
)

// nodeToOptionalObjectMeta decodes metadata after duplicate keys were rejected.
//
// api/meta owns metadata scalar and collection JSON rules. codecjson only makes
// sure the metadata node came from the ordered duplicate-checking parser before
// delegating to those existing contracts.
func nodeToOptionalObjectMeta(path jsonPath, node jsonNode, config resolvedDecodeConfig) (meta.ObjectMeta, error) {
	metadataNode, ok := node.member(apidocument.ObjectFieldMetadata.String())
	if !ok {
		return meta.ObjectMeta{}, nil
	}
	if err := requireObject(path.Member(apidocument.ObjectFieldMetadata.String()), metadataNode, "metadata must be a JSON object"); err != nil {
		return meta.ObjectMeta{}, err
	}

	data, err := jsonNodeBytes(metadataNode, false)
	if err != nil {
		return meta.ObjectMeta{}, wrapAt(
			path.Member(apidocument.ObjectFieldMetadata.String()),
			ErrInvalidEnvelope,
			codec.ErrDecodeFailed,
			ErrorReasonInvalidEnvelope,
			"metadata cannot be re-encoded for metadata decoding",
			err,
		)
	}

	var objectMeta meta.ObjectMeta
	if err := json.Unmarshal(data, &objectMeta); err != nil {
		return meta.ObjectMeta{}, wrapAt(
			path.Member(apidocument.ObjectFieldMetadata.String()),
			ErrInvalidEnvelope,
			codec.ErrInvalidDocument,
			ErrorReasonInvalidEnvelope,
			"metadata is invalid",
			err,
		)
	}

	return objectMeta, nil
}

// objectMetaToNode delegates metadata JSON shape to api/meta and re-parses it.
//
// Re-parsing keeps the final encoder on the ordered jsonNode path, so envelope
// field order and Pretty/EscapeHTML behavior remain consistent.
func objectMetaToNode(path jsonPath, objectMeta meta.ObjectMeta, config resolvedDecodeConfig) (jsonNode, error) {
	data, err := json.Marshal(objectMeta)
	if err != nil {
		return jsonNode{}, wrapAt(
			path,
			ErrInvalidEnvelope,
			errors.Join(codec.ErrEncodeFailed, codec.ErrInvalidDocument),
			ErrorReasonInvalidEnvelope,
			"metadata cannot be encoded",
			err,
		)
	}

	node, err := decodeJSONDocument(bytes.NewReader(data), config)
	if err != nil {
		return jsonNode{}, wrapAt(
			path,
			ErrInvalidEnvelope,
			errors.Join(codec.ErrEncodeFailed, codec.ErrInvalidDocument),
			ErrorReasonInvalidEnvelope,
			"metadata encoder produced invalid JSON",
			err,
		)
	}
	if err := requireObject(path, node, "metadata must encode as a JSON object"); err != nil {
		return jsonNode{}, err
	}

	return node, nil
}
