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

package types

// fieldState stores builder metadata shared by all field wrappers.
//
// It remains private so only package-owned field builders can construct
// FieldDescriptor values. The final FieldDescriptor is created by pairing this
// state with a typed builder descriptor.
type fieldState struct {
	// name is copied from the FieldBuilder entrypoint.
	name FieldName
	// presence is set by Required or Optional.
	presence FieldPresence
	// description is optional human-facing descriptor text.
	description string
}

// withRequired returns a copy marked as required.
func (s fieldState) withRequired() fieldState {
	s.presence = PresenceRequired

	return s
}

// withOptional returns a copy marked as optional.
func (s fieldState) withOptional() fieldState {
	s.presence = PresenceOptional

	return s
}

// withDescription returns a copy with human-facing descriptor text.
func (s fieldState) withDescription(text string) fieldState {
	s.description = text

	return s
}

// fieldWithType finalizes shared state with a detached value descriptor.
func (s fieldState) fieldWithType(desc Descriptor) FieldDescriptor {
	return FieldDescriptor{
		name:        s.name,
		presence:    s.presence,
		descriptor:  cloneDescriptor(desc),
		description: s.description,
	}
}
