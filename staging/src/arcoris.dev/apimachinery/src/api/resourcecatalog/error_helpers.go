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

package resourcecatalog

import (
	"fmt"

	"arcoris.dev/apimachinery/api/internal/diagnostic"
)

// catalogError creates a catalog diagnostic with structured context.
func catalogError(path string, err error, reason ErrorReason, detail string) error {
	return &Error{
		Record: diagnostic.NewRecord(path, err, reason, detail),
	}
}

// catalogErrorf formats a catalog diagnostic detail.
func catalogErrorf(path string, err error, reason ErrorReason, format string, args ...any) error {
	return catalogError(path, err, reason, fmt.Sprintf(format, args...))
}

// nestedCatalogError preserves a lower-level resource diagnostic under a
// catalog registration failure.
func nestedCatalogError(path string, reason ErrorReason, detail string, cause error) error {
	return &Error{
		Record: diagnostic.WrapRecord(path, ErrInvalidCatalog, reason, detail, cause),
	}
}

// nilCatalogError creates the standard nil receiver diagnostic for write
// operations.
func nilCatalogError() error {
	return catalogError(
		"catalog",
		ErrNilCatalog,
		ErrorReasonNilCatalog,
		"catalog receiver must be non-nil",
	)
}
