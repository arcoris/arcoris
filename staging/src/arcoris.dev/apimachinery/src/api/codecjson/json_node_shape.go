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

import "arcoris.dev/apimachinery/api/codec"

// fieldAllowed reports whether a JSON object member belongs to a fixed shape.
type fieldAllowed func(name string) bool

// requireObject verifies that a node is a JSON object at a typed boundary.
//
// The generic value codec accepts any JSON value, but object envelopes and
// ownership state have fixed object shapes. Keeping this helper shared makes
// those diagnostics consistent without introducing a public schema layer.
func requireObject(path jsonPath, node jsonNode, detail string) error {
	if node.kind == jsonKindObject {
		return nil
	}

	return errorAt(path, ErrInvalidEnvelope, codec.ErrInvalidDocument, ErrorReasonInvalidEnvelope, detail)
}

// expectString extracts a JSON string from a typed envelope/state field.
func expectString(path jsonPath, node jsonNode, detail string) (string, error) {
	if node.kind == jsonKindString {
		return node.stringValue, nil
	}

	return "", errorAt(path, ErrInvalidEnvelope, codec.ErrInvalidDocument, ErrorReasonInvalidEnvelope, detail)
}

// rejectUnknownMembers rejects fields outside one fixed JSON document shape.
//
// A predicate avoids package-level mutable maps while still making allowed
// fields explicit beside each concrete document model.
func rejectUnknownMembers(path jsonPath, node jsonNode, allow fieldAllowed, detailName string) error {
	for _, member := range node.members {
		if allow(member.name) {
			continue
		}

		return errorfAt(
			path.Member(member.name),
			ErrInvalidEnvelope,
			codec.ErrInvalidDocument,
			ErrorReasonInvalidEnvelope,
			"unknown %s field %q",
			detailName,
			member.name,
		)
	}

	return nil
}
