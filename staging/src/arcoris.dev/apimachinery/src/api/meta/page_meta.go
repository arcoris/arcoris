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

// PageMeta represents metadata for paginated collection responses.
//
// It does not contain items, selectors, query filters, storage cursor
// internals, watch state, or route information.
type PageMeta struct {
	// ResourceVersion is the opaque consistency/change token for the collection page.
	ResourceVersion stamp.ResourceVersion `json:"resourceVersion,omitempty"`
	// ContinueToken is the opaque pagination token encoded as "continue" on the wire.
	ContinueToken PageToken `json:"continue,omitempty"`
	// RemainingItemCount is an optional server-provided count hint.
	RemainingItemCount *uint64 `json:"remainingItemCount,omitempty"`
}

// IsZero reports whether all page metadata is absent.
func (m PageMeta) IsZero() bool {
	return m.ResourceVersion.IsZero() &&
		m.ContinueToken.IsZero() &&
		m.RemainingItemCount == nil
}
