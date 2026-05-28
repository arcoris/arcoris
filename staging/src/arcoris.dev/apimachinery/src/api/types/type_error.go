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
	"fmt"
)

var (
	// ErrInvalidType classifies malformed structural type descriptors.
	ErrInvalidType = errors.New("invalid type")
	// ErrInvalidTypeCode classifies unsupported TypeCode values.
	ErrInvalidTypeCode = errors.New("invalid type code")
	// ErrInvalidField classifies malformed object field descriptors.
	ErrInvalidField = errors.New("invalid field")
	// ErrDuplicateField classifies repeated field names within one object.
	ErrDuplicateField = errors.New("duplicate field")
	// ErrInvalidTypeReference classifies malformed TypeRef names or ref cycles.
	ErrInvalidTypeReference = errors.New("invalid type reference")
	// ErrUnknownTypeReference classifies TypeRef names absent from a Resolver.
	ErrUnknownTypeReference = errors.New("unknown type reference")
)

// TypeError attaches a descriptor path to a classified validation error.
//
// Path is a descriptor path such as object.fields[spec].type, list.elem, or
// ref(arcoris.meta.Name). It is not a path into a future concrete API object.
type TypeError struct {
	// Path identifies the descriptor location that failed validation.
	Path string
	// Err is the classified validation error.
	Err error
}

// Error returns a stable diagnostic message for e.
func (e *TypeError) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.Path == "" {
		return fmt.Sprintf("types: %v", e.Err)
	}
	return fmt.Sprintf("types: %s: %v", e.Path, e.Err)
}

// Unwrap returns the classified validation error.
func (e *TypeError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

// typeError creates a path-aware validation error.
func typeError(path string, err error) error {
	return &TypeError{Path: path, Err: err}
}
