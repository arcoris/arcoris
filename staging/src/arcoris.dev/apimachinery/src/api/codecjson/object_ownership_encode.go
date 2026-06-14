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
// Encoding validates objectownership.State before writing. The codec does not
// normalize, compute ownership, or interpret lifecycle semantics; ownership
// canonicalization belongs to objectownership constructors and objectstore
// commit boundaries.
func ownershipStateToNode(path jsonPath, state objectownership.State, config resolvedEncodeConfig) (jsonNode, error) {
	fields := apidocument.Fields().Ownership()

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

	metadata := state.Metadata()
	members := []jsonMember{}
	if config.emptySurfaces == jsonconfig.EmptyOwnershipSurfaceEmit || !state.Desired().IsEmpty() {
		members = append(members, jsonMember{
			name:  fields.Desired().String(),
			value: ownershipSurfaceToNode(path.Member(fields.Desired().String()), state.Desired()),
		})
	}
	if config.emptySurfaces == jsonconfig.EmptyOwnershipSurfaceEmit || !state.Observed().IsEmpty() {
		members = append(members, jsonMember{
			name:  fields.Observed().String(),
			value: ownershipSurfaceToNode(path.Member(fields.Observed().String()), state.Observed()),
		})
	}
	if config.emptySurfaces == jsonconfig.EmptyOwnershipSurfaceEmit || !metadata.IsEmpty() {
		members = append(members, jsonMember{
			name:  fields.MetadataField().String(),
			value: ownershipMetadataToNode(path.Member(fields.MetadataField().String()), metadata, config),
		})
	}

	return jsonNode{kind: jsonKindObject, members: members}, nil
}

// ownershipSurfaceToNode converts one ownership surface into an entries object.
func ownershipSurfaceToNode(path jsonPath, surface fieldownership.State) jsonNode {
	fields := apidocument.Fields().Ownership().Surface()
	entries := surface.Entries()
	items := make([]jsonNode, 0, len(entries))
	for i, entry := range entries {
		items = append(items, ownershipEntryToNode(path.Member(fields.Entries().String()).Index(i), entry))
	}

	return jsonNode{
		kind: jsonKindObject,
		members: []jsonMember{
			{
				name:  fields.Entries().String(),
				value: jsonNode{kind: jsonKindArray, items: items},
			},
		},
	}
}

// ownershipMetadataToNode converts metadata ownership surfaces.
func ownershipMetadataToNode(path jsonPath, metadata objectownership.MetadataState, config resolvedEncodeConfig) jsonNode {
	fields := apidocument.Fields().Ownership().Metadata()
	members := []jsonMember{}
	if config.emptySurfaces == jsonconfig.EmptyOwnershipSurfaceEmit || !metadata.Labels().IsEmpty() {
		members = append(members, jsonMember{
			name:  fields.Labels().String(),
			value: ownershipSurfaceToNode(path.Member(fields.Labels().String()), metadata.Labels()),
		})
	}
	if config.emptySurfaces == jsonconfig.EmptyOwnershipSurfaceEmit || !metadata.Annotations().IsEmpty() {
		members = append(members, jsonMember{
			name:  fields.Annotations().String(),
			value: ownershipSurfaceToNode(path.Member(fields.Annotations().String()), metadata.Annotations()),
		})
	}

	return jsonNode{kind: jsonKindObject, members: members}
}

// ownershipEntryToNode converts one owner entry.
func ownershipEntryToNode(path jsonPath, entry fieldownership.Entry) jsonNode {
	entryFields := apidocument.Fields().Ownership().Entry()
	paths := entry.Fields().Paths()
	fields := make([]jsonNode, 0, len(paths))
	for _, field := range paths {
		fields = append(fields, jsonNode{kind: jsonKindString, stringValue: field.String()})
	}

	return jsonNode{
		kind: jsonKindObject,
		members: []jsonMember{
			{name: entryFields.Owner().String(), value: jsonNode{kind: jsonKindString, stringValue: entry.Owner().String()}},
			{name: entryFields.Fields().String(), value: jsonNode{kind: jsonKindArray, items: fields}},
		},
	}
}
