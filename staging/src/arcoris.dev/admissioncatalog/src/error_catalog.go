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

package admissioncatalog

import (
	"errors"
	"fmt"
)

var (
	// ErrNilCatalog identifies composition with a nil catalog pointer.
	ErrNilCatalog = errors.New("admissioncatalog: nil catalog")
)

// NilCatalogError reports a nil catalog in composition input.
type NilCatalogError struct {
	// Operation names the composition operation that rejected the nil catalog.
	Operation string

	// Index identifies the nil catalog position when the operation receives a
	// slice of catalogs. A negative value means no index applies.
	Index int
}

// Error returns a concise diagnostic for the nil catalog input.
func (e NilCatalogError) Error() string {
	if e.Index >= 0 {
		return fmt.Sprintf("admissioncatalog: %s catalog[%d]: nil catalog", e.Operation, e.Index)
	}
	if e.Operation != "" {
		return fmt.Sprintf("admissioncatalog: %s: nil catalog", e.Operation)
	}
	return ErrNilCatalog.Error()
}

// Unwrap returns the sentinel error for errors.Is.
func (e NilCatalogError) Unwrap() error {
	return ErrNilCatalog
}
