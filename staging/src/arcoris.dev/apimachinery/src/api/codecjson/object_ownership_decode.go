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
	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/objectownership"
)

// nodeToOwnershipState converts lower-camel JSON into canonical ownership state.
//
// Decoding targets objectownership.State, not a raw JSON document model.
// Duplicate owner entries and duplicate fields may collapse through
// fieldownership.State construction. That loss of raw duplication is
// intentional because the codec returns the semantic ownership model.
func nodeToOwnershipState(path jsonPath, node jsonNode, config resolvedDecodeConfig) (objectownership.State, error) {
	if err := requireObject(path, node, "object ownership root must be a JSON object"); err != nil {
		return objectownership.State{}, err
	}
	if config.rejectUnknownOwnershipFields {
		if err := rejectUnknownMembers(path, node, allowOwnershipStateField, "object ownership state"); err != nil {
			return objectownership.State{}, err
		}
	}

	desired := fieldownership.EmptyState()
	if desiredNode, ok := node.member(apidocument.OwnershipFieldDesired.String()); ok {
		surface, err := nodeToOwnershipSurface(path.Member(apidocument.OwnershipFieldDesired.String()), desiredNode, config)
		if err != nil {
			return objectownership.State{}, err
		}
		desired = surface
	}

	observed := fieldownership.EmptyState()
	if observedNode, ok := node.member(apidocument.OwnershipFieldObserved.String()); ok {
		surface, err := nodeToOwnershipSurface(path.Member(apidocument.OwnershipFieldObserved.String()), observedNode, config)
		if err != nil {
			return objectownership.State{}, err
		}
		observed = surface
	}

	metadata := objectownership.NewMetadataState(fieldownership.EmptyState(), fieldownership.EmptyState())
	if metadataNode, ok := node.member(apidocument.OwnershipFieldMetadata.String()); ok {
		decoded, err := nodeToOwnershipMetadata(path.Member(apidocument.OwnershipFieldMetadata.String()), metadataNode, config)
		if err != nil {
			return objectownership.State{}, err
		}
		metadata = decoded
	}

	state := objectownership.NewStateWithSurfaces(desired, observed, metadata)
	if config.validateOwnershipState {
		if err := objectownership.Validate(state); err != nil {
			return objectownership.State{}, wrapAt(
				path,
				ErrInvalidEnvelope,
				codec.ErrInvalidDocument,
				ErrorReasonInvalidEnvelope,
				"object ownership state is invalid",
				err,
			)
		}
	}

	return state, nil
}

// nodeToOwnershipMetadata converts lower-camel metadata ownership object.
func nodeToOwnershipMetadata(path jsonPath, node jsonNode, config resolvedDecodeConfig) (objectownership.MetadataState, error) {
	if err := requireObject(path, node, "object ownership metadata must be a JSON object"); err != nil {
		return objectownership.MetadataState{}, err
	}
	if config.rejectUnknownOwnershipFields {
		if err := rejectUnknownMembers(path, node, allowOwnershipMetadataField, "object ownership metadata"); err != nil {
			return objectownership.MetadataState{}, err
		}
	}

	labels := fieldownership.EmptyState()
	if labelsNode, ok := node.member(apidocument.OwnershipFieldLabels.String()); ok {
		surface, err := nodeToOwnershipSurface(path.Member(apidocument.OwnershipFieldLabels.String()), labelsNode, config)
		if err != nil {
			return objectownership.MetadataState{}, err
		}
		labels = surface
	}

	annotations := fieldownership.EmptyState()
	if annotationsNode, ok := node.member(apidocument.OwnershipFieldAnnotations.String()); ok {
		surface, err := nodeToOwnershipSurface(path.Member(apidocument.OwnershipFieldAnnotations.String()), annotationsNode, config)
		if err != nil {
			return objectownership.MetadataState{}, err
		}
		annotations = surface
	}

	return objectownership.NewMetadataState(labels, annotations), nil
}

// nodeToOwnershipSurface converts one surface object into field ownership state.
func nodeToOwnershipSurface(path jsonPath, node jsonNode, config resolvedDecodeConfig) (fieldownership.State, error) {
	if err := requireObject(path, node, "object ownership surface must be a JSON object"); err != nil {
		return fieldownership.State{}, err
	}
	if config.rejectUnknownOwnershipFields {
		if err := rejectUnknownMembers(path, node, allowOwnershipSurfaceField, "object ownership surface"); err != nil {
			return fieldownership.State{}, err
		}
	}

	entriesNode, ok := node.member(apidocument.OwnershipFieldEntries.String())
	if !ok {
		return fieldownership.EmptyState(), nil
	}
	if entriesNode.kind != jsonKindArray {
		return fieldownership.State{}, errorAt(
			path.Member(apidocument.OwnershipFieldEntries.String()),
			ErrInvalidEnvelope,
			codec.ErrInvalidDocument,
			ErrorReasonInvalidEnvelope,
			"object ownership entries must be a JSON array",
		)
	}

	entries := make([]fieldownership.Entry, 0, len(entriesNode.items))
	for i, item := range entriesNode.items {
		entry, err := nodeToOwnershipEntry(path.Member(apidocument.OwnershipFieldEntries.String()).Index(i), item, config)
		if err != nil {
			return fieldownership.State{}, err
		}
		entries = append(entries, entry)
	}

	state, err := fieldownership.NewState(entries...)
	if err != nil {
		return fieldownership.State{}, wrapAt(
			path,
			ErrInvalidEnvelope,
			codec.ErrInvalidDocument,
			ErrorReasonInvalidEnvelope,
			"object ownership surface is invalid",
			err,
		)
	}

	return state, nil
}

// nodeToOwnershipEntry converts one ownership entry object.
func nodeToOwnershipEntry(path jsonPath, node jsonNode, config resolvedDecodeConfig) (fieldownership.Entry, error) {
	if err := requireObject(path, node, "object ownership entry must be a JSON object"); err != nil {
		return fieldownership.Entry{}, err
	}
	if config.rejectUnknownOwnershipFields {
		if err := rejectUnknownMembers(path, node, allowOwnershipEntryField, "object ownership entry"); err != nil {
			return fieldownership.Entry{}, err
		}
	}

	owner, err := nodeToRequiredOwnershipOwner(path, node)
	if err != nil {
		return fieldownership.Entry{}, err
	}
	fields, err := nodeToOwnershipFields(path, node)
	if err != nil {
		return fieldownership.Entry{}, err
	}

	entry, err := fieldownership.NewEntry(owner, fields)
	if err != nil {
		return fieldownership.Entry{}, wrapAt(
			path,
			ErrInvalidEnvelope,
			codec.ErrInvalidDocument,
			ErrorReasonInvalidEnvelope,
			"object ownership entry is invalid",
			err,
		)
	}

	return entry, nil
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

// nodeToOwnershipFields extracts an optional canonical field path array.
func nodeToOwnershipFields(path jsonPath, node jsonNode) (fieldpath.Set, error) {
	fieldsNode, ok := node.member(apidocument.OwnershipFieldFields.String())
	if !ok {
		return fieldpath.EmptySet(), nil
	}
	if fieldsNode.kind != jsonKindArray {
		return fieldpath.Set{}, errorAt(
			path.Member(apidocument.OwnershipFieldFields.String()),
			ErrInvalidEnvelope,
			codec.ErrInvalidDocument,
			ErrorReasonInvalidEnvelope,
			"object ownership entry fields must be a JSON array",
		)
	}

	paths := make([]fieldpath.Path, 0, len(fieldsNode.items))
	for i, item := range fieldsNode.items {
		field, err := expectString(
			path.Member(apidocument.OwnershipFieldFields.String()).Index(i),
			item,
			"object ownership entry field must be a JSON string",
		)
		if err != nil {
			return fieldpath.Set{}, err
		}
		parsed, err := fieldpath.ParseCanonical(field)
		if err != nil {
			return fieldpath.Set{}, wrapAt(
				path.Member(apidocument.OwnershipFieldFields.String()).Index(i),
				ErrInvalidEnvelope,
				codec.ErrInvalidDocument,
				ErrorReasonInvalidEnvelope,
				"object ownership entry field path is invalid",
				err,
			)
		}
		paths = append(paths, parsed)
	}

	fields, err := fieldpath.NewSet(paths...)
	if err != nil {
		return fieldpath.Set{}, wrapAt(
			path.Member(apidocument.OwnershipFieldFields.String()),
			ErrInvalidEnvelope,
			codec.ErrInvalidDocument,
			ErrorReasonInvalidEnvelope,
			"object ownership entry fields are invalid",
			err,
		)
	}

	return fields, nil
}
