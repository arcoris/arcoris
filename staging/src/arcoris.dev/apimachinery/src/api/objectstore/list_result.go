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

package objectstore

// ListResult is the successful result of a collection read.
type ListResult struct {
	// Items are detached live states that matched the request.
	Items []ListItem

	// Revision is the store revision observed by the list operation.
	Revision Revision
}

// IsZero reports whether r contains no items and no observed revision.
func (r ListResult) IsZero() bool {
	return len(r.Items) == 0 && r.Revision.IsZero()
}

// Len returns the number of listed live items.
func (r ListResult) Len() int {
	return len(r.Items)
}

// Clone returns a detached copy of r.
func (r ListResult) Clone() ListResult {
	if r.Items == nil {
		return ListResult{Revision: r.Revision}
	}

	items := make([]ListItem, len(r.Items))
	for i, item := range r.Items {
		items[i] = item.Clone()
	}

	return ListResult{Items: items, Revision: r.Revision}
}
