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
	"arcoris.dev/apimachinery/api/fieldownership"
	"arcoris.dev/apimachinery/api/objectownership"
)

// ownershipStateToNode converts canonical ownership state into lower-camel JSON.
//
// Encoding validates the state before writing. Deterministic encoding may run
// objectownership.Normalize, which is intentionally idempotent for public
// objectownership.State values.
func ownershipStateToNode(path jsonPath, state objectownership.State, config resolvedEncodeConfig) (jsonNode, error) {
	if err := objectownership.Validate(state); err != nil {
		return jsonNode{}, wrapAt(
			path,
			ErrInvalidEnvelope,
			errors.Join(codec.ErrEncodeFailed, codec.ErrInvalidDocument),
			ErrorReasonInvalidEnvelope,
			"object ownership state is invalid",
			err,
		)
	}
	if shouldNormalizeOwnership(config) {
		state = objectownership.Normalize(state)
	}

	metadata := state.Metadata()
	members := []jsonMember{}
	if config.emptySurfaces == jsonconfig.EmptyOwnershipSurfaceEmit || !state.Desired().IsEmpty() {
		members = append(members, jsonMember{
			name:  apidocument.OwnershipFieldDesired.String(),
			value: ownershipSurfaceToNode(path.Member(apidocument.OwnershipFieldDesired.String()), state.Desired(), config),
		})
	}
	if config.emptySurfaces == jsonconfig.EmptyOwnershipSurfaceEmit || !state.Observed().IsEmpty() {
		members = append(members, jsonMember{
			name:  apidocument.OwnershipFieldObserved.String(),
			value: ownershipSurfaceToNode(path.Member(apidocument.OwnershipFieldObserved.String()), state.Observed(), config),
		})
	}
	if config.emptySurfaces == jsonconfig.EmptyOwnershipSurfaceEmit || !metadata.IsEmpty() {
		members = append(members, jsonMember{
			name:  apidocument.OwnershipFieldMetadata.String(),
			value: ownershipMetadataToNode(path.Member(apidocument.OwnershipFieldMetadata.String()), metadata, config),
		})
	}

	return jsonNode{kind: jsonKindObject, members: members}, nil
}

// shouldNormalizeOwnership reports whether config asks for canonical ownership order.
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

// ownershipSurfaceToNode converts one ownership surface into an entries object.
func ownershipSurfaceToNode(path jsonPath, surface fieldownership.State, config resolvedEncodeConfig) jsonNode {
	entries := surface.Entries()
	items := make([]jsonNode, 0, len(entries))
	for i, entry := range entries {
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

// ownershipMetadataToNode converts metadata ownership surfaces.
func ownershipMetadataToNode(path jsonPath, metadata objectownership.MetadataState, config resolvedEncodeConfig) jsonNode {
	members := []jsonMember{}
	if config.emptySurfaces == jsonconfig.EmptyOwnershipSurfaceEmit || !metadata.Labels().IsEmpty() {
		members = append(members, jsonMember{
			name:  apidocument.OwnershipFieldLabels.String(),
			value: ownershipSurfaceToNode(path.Member(apidocument.OwnershipFieldLabels.String()), metadata.Labels(), config),
		})
	}
	if config.emptySurfaces == jsonconfig.EmptyOwnershipSurfaceEmit || !metadata.Annotations().IsEmpty() {
		members = append(members, jsonMember{
			name:  apidocument.OwnershipFieldAnnotations.String(),
			value: ownershipSurfaceToNode(path.Member(apidocument.OwnershipFieldAnnotations.String()), metadata.Annotations(), config),
		})
	}

	return jsonNode{kind: jsonKindObject, members: members}
}

// ownershipEntryToNode converts one owner entry.
func ownershipEntryToNode(path jsonPath, entry fieldownership.Entry) jsonNode {
	paths := entry.Fields().Paths()
	fields := make([]jsonNode, 0, len(paths))
	for _, field := range paths {
		fields = append(fields, jsonNode{kind: jsonKindString, stringValue: field.String()})
	}

	return jsonNode{
		kind: jsonKindObject,
		members: []jsonMember{
			{name: apidocument.OwnershipFieldOwner.String(), value: jsonNode{kind: jsonKindString, stringValue: entry.Owner().String()}},
			{name: apidocument.OwnershipFieldFields.String(), value: jsonNode{kind: jsonKindArray, items: fields}},
		},
	}
}
