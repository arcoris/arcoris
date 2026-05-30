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

package meta

import "arcoris.dev/apimachinery/api/meta/stamp"

// ListMeta represents metadata for list and page responses.
//
// It does not contain list items, selectors, query filters, storage cursor
// internals, watch state, or route information.
type ListMeta struct {
	// ResourceVersion is the opaque consistency/change token for the list response.
	ResourceVersion stamp.ResourceVersion
	// ContinueToken is an opaque pagination token supplied by a higher layer.
	ContinueToken PageToken
	// RemainingItemCount is an optional server-provided count hint.
	RemainingItemCount *uint64
}

// IsZero reports whether all list metadata is absent.
func (m ListMeta) IsZero() bool {
	return m.ResourceVersion.IsZero() &&
		m.ContinueToken.IsZero() &&
		m.RemainingItemCount == nil
}
