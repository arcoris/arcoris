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

package objectwatch

import "arcoris.dev/apimachinery/api/objectstore"

// Request identifies one structural object collection watch.
//
// Request intentionally contains no objectquery.Query. objectwatch defines the
// committed structural stream only; filtered watch semantics are built by
// composing EventChanged.Change with objectquery outside this package.
type Request struct {
	// Collection is the structural resource/scope collection to watch. It
	// follows the same validation rules as objectstore.ListRequest, but
	// objectwatch never executes a list operation.
	Collection objectstore.ListRequest
	// Start defines whether the stream catches up after a revision or starts at
	// the source's current progress point.
	Start Start
	// AllowProgress permits, but does not require, EventProgress markers.
	AllowProgress bool
}

// IsValid reports whether r passes request validation.
func (r Request) IsValid() bool {
	return r.Validate() == nil
}
