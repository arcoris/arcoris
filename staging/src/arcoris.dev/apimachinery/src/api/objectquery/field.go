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

package objectquery

import (
	"fmt"

	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/objectsurface"
)

// FieldRef identifies one registered selectable field on an object surface.
//
// A FieldRef is not an arbitrary JSONPath. It becomes queryable only when a
// SelectableFieldSet explicitly resolves it during compilation.
type FieldRef struct {
	// Surface is the semantic object surface that owns the field.
	Surface objectsurface.Kind
	// Path is a surface-relative semantic payload path, not JSONPath.
	Path fieldpath.Path
}

// String returns a stable diagnostic field reference.
func (r FieldRef) String() string {
	if r.Path.IsRoot() {
		return r.Surface.String()
	}

	return r.Surface.String() + r.Path.CanonicalText()[1:]
}

// Validate checks the field reference shape before registry resolution.
func (r FieldRef) Validate() error {
	if !r.Surface.IsValid() {
		return fmt.Errorf("%w: surface %q", ErrInvalidField, r.Surface.String())
	}
	if !isFieldSurfaceSupported(r.Surface) {
		return fmt.Errorf("%w: unsupported field surface %q", ErrInvalidField, r.Surface.String())
	}
	if r.Path.IsRoot() {
		return fmt.Errorf("%w: empty field path", ErrInvalidField)
	}
	if err := r.Path.ValidateStructure(); err != nil {
		return fmt.Errorf("%w: invalid field path: %w", ErrInvalidField, err)
	}

	return nil
}

// isFieldSurfaceSupported reports whether objectquery can evaluate field
// payloads on surface. Metadata surfaces use dedicated metadata terms instead.
func isFieldSurfaceSupported(surface objectsurface.Kind) bool {
	kinds := objectsurface.Kinds()
	return surface == kinds.Desired() || surface == kinds.Observed()
}

// fieldRefKey is the private structural key used by StaticFieldSet.
type fieldRefKey struct {
	surface objectsurface.Kind
	path    string
}

// key returns a stable semantic key without using FieldRef.String diagnostics.
func (r FieldRef) key() fieldRefKey {
	return fieldRefKey{surface: r.Surface, path: r.Path.CanonicalText()}
}
