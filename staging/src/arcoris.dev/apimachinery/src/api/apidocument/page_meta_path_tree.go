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

package apidocument

// PageMetaPaths groups paths for list/page metadata documents.
type PageMetaPaths struct {
	// base is the document path root for this PageMeta view.
	base Path
}

// ResourceVersion returns the page metadata resourceVersion path.
func (p PageMetaPaths) ResourceVersion() Path {
	return p.base.Child(Fields().PageMeta().ResourceVersion())
}

// Continue returns the page metadata continue path.
func (p PageMetaPaths) Continue() Path {
	return p.base.Child(Fields().PageMeta().Continue())
}

// RemainingItemCount returns the page metadata remainingItemCount path.
func (p PageMetaPaths) RemainingItemCount() Path {
	return p.base.Child(Fields().PageMeta().RemainingItemCount())
}
