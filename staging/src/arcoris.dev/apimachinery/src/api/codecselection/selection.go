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

package codecselection

import (
	"arcoris.dev/apimachinery/api/codec"
	"arcoris.dev/apimachinery/api/codecregistry"
)

// Selection describes the exact configured codec candidate selected at runtime.
type Selection struct {
	// Direction records whether the selected binding is for decode or encode.
	Direction Direction

	// Transport records the selected byte or stream transport.
	Transport Transport

	// Target records the selected API document model.
	Target codec.Target

	// ContentType records the normalized content key that matched.
	ContentType ContentType

	// EntryID records the exact configured codec candidate identity.
	EntryID codecregistry.EntryID

	// Entry records the registry metadata and implementation for EntryID.
	//
	// Selection exposes Entry for inspection and advanced diagnostics. Typed
	// Select* methods already return the selected capability directly, so most
	// callers should use that typed value instead of asserting Entry.Codec().
	Entry codecregistry.Entry
}

// IsZero reports whether s contains no selected entry.
func (s Selection) IsZero() bool {
	return s.Direction.IsZero() &&
		s.Transport.IsZero() &&
		s.Target.IsZero() &&
		s.ContentType.IsZero() &&
		s.EntryID.IsZero() &&
		s.Entry.IsZero()
}
