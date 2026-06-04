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
	"arcoris.dev/apimachinery/api/objectownership"
)

// ownershipDocumentToNode converts an ownership document into lower-camel JSON.
//
// Encoding validates the document before writing. Deterministic encoding uses
// objectownership.Normalize so owner and field ordering match the semantic
// ownership package rather than JSON-specific sorting rules.
func ownershipDocumentToNode(path jsonPath, doc objectownership.Document, config resolvedEncodeConfig) (jsonNode, error) {
	if err := objectownership.Validate(doc); err != nil {
		return jsonNode{}, wrapAt(
			path,
			ErrInvalidEnvelope,
			errors.Join(codec.ErrEncodeFailed, codec.ErrInvalidDocument),
			ErrorReasonInvalidEnvelope,
			"object ownership document is invalid",
			err,
		)
	}
	if shouldNormalizeOwnership(config) {
		normalized, err := objectownership.Normalize(doc)
		if err != nil {
			return jsonNode{}, wrapAt(
				path,
				ErrInvalidEnvelope,
				errors.Join(codec.ErrEncodeFailed, codec.ErrInvalidDocument),
				ErrorReasonInvalidEnvelope,
				"object ownership document cannot be normalized",
				err,
			)
		}
		doc = normalized
	}

	members := []jsonMember{
		{name: apidocument.OwnershipFieldVersion.String(), value: jsonNode{kind: jsonKindString, stringValue: doc.Version.String()}},
	}
	if config.emptyDesired == jsonconfig.EmptyOwnershipSurfaceEmit || !doc.Desired.IsEmpty() {
		members = append(members, jsonMember{
			name:  apidocument.OwnershipFieldDesired.String(),
			value: ownershipSurfaceToNode(path.Member(apidocument.OwnershipFieldDesired.String()), doc.Desired, config),
		})
	}

	return jsonNode{kind: jsonKindObject, members: members}, nil
}

// shouldNormalizeOwnership reports whether config asks for canonical document order.
func shouldNormalizeOwnership(config resolvedEncodeConfig) bool {
	switch config.ownershipNormalize {
	case jsonconfig.OwnershipNormalizeAlways:
		return true
	case jsonconfig.OwnershipNormalizeWhenDeterministic:
		return config.deterministic
	default:
		return false
	}
}

// ownershipSurfaceToNode converts one surface while preserving entry order.
func ownershipSurfaceToNode(path jsonPath, surface objectownership.Surface, config resolvedEncodeConfig) jsonNode {
	items := make([]jsonNode, 0, len(surface.Entries))
	for i, entry := range surface.Entries {
		items = append(items, ownershipEntryToNode(path.Member(apidocument.OwnershipFieldEntries.String()).Index(i), entry))
	}

	members := []jsonMember{}
	if config.emptyEntries == jsonconfig.EmptyEntriesEmit || len(items) > 0 {
		members = append(members, jsonMember{
			name:  apidocument.OwnershipFieldEntries.String(),
			value: jsonNode{kind: jsonKindArray, items: items},
		})
	}

	return jsonNode{kind: jsonKindObject, members: members}
}

// ownershipEntryToNode converts one entry while preserving field order.
func ownershipEntryToNode(path jsonPath, entry objectownership.Entry) jsonNode {
	fields := make([]jsonNode, 0, len(entry.Fields))
	for _, field := range entry.Fields {
		fields = append(fields, jsonNode{kind: jsonKindString, stringValue: field.String()})
	}

	return jsonNode{
		kind: jsonKindObject,
		members: []jsonMember{
			{name: apidocument.OwnershipFieldOwner.String(), value: jsonNode{kind: jsonKindString, stringValue: entry.Owner.String()}},
			{name: apidocument.OwnershipFieldFields.String(), value: jsonNode{kind: jsonKindArray, items: fields}},
		},
	}
}
