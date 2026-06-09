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

import (
	"errors"

	"arcoris.dev/apimachinery/api/internal/diagnostic"
)

var (
	// ErrInvalidDescriptor classifies malformed structural descriptors.
	ErrInvalidDescriptor = errors.New("invalid descriptor")
	// ErrInvalidDescriptorKind classifies unsupported DescriptorKind values.
	ErrInvalidDescriptorKind = errors.New("invalid descriptor kind")
	// ErrInvalidField classifies malformed object field descriptors.
	ErrInvalidField = errors.New("invalid field")
	// ErrDuplicateField classifies repeated field names within one object.
	ErrDuplicateField = errors.New("duplicate field")
	// ErrInvalidDescriptorReference classifies malformed DescriptorRef names or ref cycles.
	ErrInvalidDescriptorReference = errors.New("invalid descriptor reference")
	// ErrUnresolvedDescriptorReference classifies DescriptorRef names absent from a Resolver.
	ErrUnresolvedDescriptorReference = errors.New("unresolved descriptor reference")
)

// DescriptorError attaches structured descriptor diagnostics to a classified error.
//
// Path is a descriptor path such as object.fields[spec].type, list.elem, or
// ref(meta.arcoris.dev.Name). It is not a path into a future concrete API object.
type DescriptorError struct {
	// Record stores the shared path, sentinel, reason, and detail fields.
	diagnostic.Record[DescriptorErrorReason]
}

// Error returns a stable diagnostic message for e.
func (e *DescriptorError) Error() string {
	if e == nil {
		return "<nil>"
	}

	return e.Record.Format("types")
}

// Unwrap returns the classified validation error.
func (e *DescriptorError) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.Record.Unwrap()
}
