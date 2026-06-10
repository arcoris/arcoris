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
	"arcoris.dev/apimachinery/api/fieldownership"
	"arcoris.dev/apimachinery/api/objectownership"
)

// nodeToOwnershipDocument converts lower-camel JSON into an ownership document.
//
// Decode preserves the raw valid document order. Canonicalization is reserved
// for deterministic encoding or objectownership.Normalize callers.
func nodeToOwnershipDocument(path jsonPath, node jsonNode, config resolvedDecodeConfig) (objectownership.Document, error) {
	if err := requireObject(path, node, "object ownership document root must be a JSON object"); err != nil {
		return objectownership.Document{}, err
	}
	if config.rejectUnknownOwnershipFields {
		if err := rejectUnknownMembers(path, node, allowOwnershipDocumentField, "object ownership document"); err != nil {
			return objectownership.Document{}, err
		}
	}

	version, err := nodeToRequiredOwnershipDocumentVersion(path, node)
	if err != nil {
		return objectownership.Document{}, err
	}

	doc := objectownership.Document{Version: version}
	if desiredNode, ok := node.member(apidocument.OwnershipFieldDesired.String()); ok {
		desired, err := nodeToOwnershipSurface(path.Member(apidocument.OwnershipFieldDesired.String()), desiredNode, config)
		if err != nil {
			return objectownership.Document{}, err
		}
		doc.Desired = desired
	}
	if config.validateOwnershipDocument {
		if err := objectownership.Validate(doc); err != nil {
			return objectownership.Document{}, wrapAt(
				path,
				ErrInvalidEnvelope,
				codec.ErrInvalidDocument,
				ErrorReasonInvalidEnvelope,
				"object ownership document is invalid",
				err,
			)
		}
	}

	return doc, nil
}

// nodeToRequiredOwnershipDocumentVersion extracts the explicit document version string.
func nodeToRequiredOwnershipDocumentVersion(path jsonPath, node jsonNode) (objectownership.DocumentVersion, error) {
	versionNode, ok := node.member(apidocument.OwnershipFieldVersion.String())
	if !ok {
		return "", errorAt(
			path.Member(apidocument.OwnershipFieldVersion.String()),
			ErrInvalidEnvelope,
			codec.ErrInvalidDocument,
			ErrorReasonInvalidEnvelope,
			"object ownership version is required",
		)
	}
	version, err := expectString(path.Member(apidocument.OwnershipFieldVersion.String()), versionNode, "object ownership version must be a JSON string")
	if err != nil {
		return "", err
	}

	return objectownership.DocumentVersion(version), nil
}

// nodeToOwnershipSurface converts a lower-camel surface object.
func nodeToOwnershipSurface(path jsonPath, node jsonNode, config resolvedDecodeConfig) (objectownership.Surface, error) {
	if err := requireObject(path, node, "object ownership desired surface must be a JSON object"); err != nil {
		return objectownership.Surface{}, err
	}
	if config.rejectUnknownOwnershipFields {
		if err := rejectUnknownMembers(path, node, allowOwnershipSurfaceField, "object ownership surface"); err != nil {
			return objectownership.Surface{}, err
		}
	}

	entriesNode, ok := node.member(apidocument.OwnershipFieldEntries.String())
	if !ok {
		return objectownership.Surface{}, nil
	}
	if entriesNode.kind != jsonKindArray {
		return objectownership.Surface{}, errorAt(
			path.Member(apidocument.OwnershipFieldEntries.String()),
			ErrInvalidEnvelope,
			codec.ErrInvalidDocument,
			ErrorReasonInvalidEnvelope,
			"object ownership entries must be a JSON array",
		)
	}

	entries := make([]objectownership.Entry, 0, len(entriesNode.items))
	for i, item := range entriesNode.items {
		entry, err := nodeToOwnershipEntry(path.Member(apidocument.OwnershipFieldEntries.String()).Index(i), item, config)
		if err != nil {
			return objectownership.Surface{}, err
		}
		entries = append(entries, entry)
	}

	return objectownership.Surface{Entries: entries}, nil
}

// nodeToOwnershipEntry converts one ownership entry object.
func nodeToOwnershipEntry(path jsonPath, node jsonNode, config resolvedDecodeConfig) (objectownership.Entry, error) {
	if err := requireObject(path, node, "object ownership entry must be a JSON object"); err != nil {
		return objectownership.Entry{}, err
	}
	if config.rejectUnknownOwnershipFields {
		if err := rejectUnknownMembers(path, node, allowOwnershipEntryField, "object ownership entry"); err != nil {
			return objectownership.Entry{}, err
		}
	}

	owner, err := nodeToRequiredOwnershipOwner(path, node)
	if err != nil {
		return objectownership.Entry{}, err
	}
	fields, err := nodeToOwnershipFields(path, node)
	if err != nil {
		return objectownership.Entry{}, err
	}

	return objectownership.Entry{
		Owner:  owner,
		Fields: fields,
	}, nil
}

// nodeToRequiredOwnershipOwner extracts the required owner identity string.
func nodeToRequiredOwnershipOwner(path jsonPath, node jsonNode) (fieldownership.Owner, error) {
	ownerNode, ok := node.member(apidocument.OwnershipFieldOwner.String())
	if !ok {
		return fieldownership.Owner{}, errorAt(
			path.Member(apidocument.OwnershipFieldOwner.String()),
			ErrInvalidEnvelope,
			codec.ErrInvalidDocument,
			ErrorReasonInvalidEnvelope,
			"object ownership entry owner is required",
		)
	}
	owner, err := expectString(path.Member(apidocument.OwnershipFieldOwner.String()), ownerNode, "object ownership entry owner must be a JSON string")
	if err != nil {
		return fieldownership.Owner{}, err
	}

	parsed, err := fieldownership.NewOwner(owner)
	if err != nil {
		return fieldownership.Owner{}, wrapAt(
			path.Member(apidocument.OwnershipFieldOwner.String()),
			ErrInvalidEnvelope,
			codec.ErrInvalidDocument,
			ErrorReasonInvalidEnvelope,
			"object ownership entry owner is invalid",
			err,
		)
	}

	return parsed, nil
}

// nodeToOwnershipFields extracts an optional fields array.
func nodeToOwnershipFields(path jsonPath, node jsonNode) ([]objectownership.Path, error) {
	fieldsNode, ok := node.member(apidocument.OwnershipFieldFields.String())
	if !ok {
		return nil, nil
	}
	if fieldsNode.kind != jsonKindArray {
		return nil, errorAt(
			path.Member(apidocument.OwnershipFieldFields.String()),
			ErrInvalidEnvelope,
			codec.ErrInvalidDocument,
			ErrorReasonInvalidEnvelope,
			"object ownership entry fields must be a JSON array",
		)
	}

	fields := make([]objectownership.Path, 0, len(fieldsNode.items))
	for i, item := range fieldsNode.items {
		field, err := expectString(
			path.Member(apidocument.OwnershipFieldFields.String()).Index(i),
			item,
			"object ownership entry field must be a JSON string",
		)
		if err != nil {
			return nil, err
		}
		fields = append(fields, objectownership.Path(field))
	}

	return fields, nil
}
