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

package objectownership

import (
	"arcoris.dev/apimachinery/api/fieldownership"
	"arcoris.dev/apimachinery/api/fieldpath"
)

// ToDocument converts in-memory ownership state into a deterministic document.
//
// The conversion is semantic only: it does not serialize, allocate storage
// metadata, or attach the document to ObjectMeta. Empty owner entries are
// omitted and all owners and paths follow fieldownership deterministic ordering.
func ToDocument(state State) Document {
	return Document{
		Version: VersionV1,
		Desired: surfaceToDocument(
			state.Desired(),
		),
	}
}

// surfaceToDocument converts one field ownership surface to document form.
//
// fieldownership.State already owns owner ordering and field set ordering, so
// this helper only filters empty entries and formats paths.
func surfaceToDocument(state fieldownership.State) Surface {
	entries := state.Entries()
	if len(entries) == 0 {
		return Surface{}
	}

	out := make([]Entry, 0, len(entries))
	for _, entry := range entries {
		if entry.IsEmpty() {
			continue
		}
		out = append(out, entryToDocument(entry))
	}

	return Surface{Entries: out}
}

// entryToDocument converts one owner entry without changing ownership semantics.
//
// Shared ownership is preserved because conversion happens per owner and does
// not inspect other owners' fields.
func entryToDocument(entry fieldownership.Entry) Entry {
	return Entry{
		Owner:  entry.Owner(),
		Fields: pathsToDocument(entry.Fields()),
	}
}

// pathsToDocument formats semantic paths in fieldpath.Set canonical order.
//
// Parent and descendant paths remain separate; the document layer does not
// compact or reinterpret fieldpath.Set semantics.
func pathsToDocument(fields fieldpath.Set) []Path {
	paths := fields.Paths()
	if len(paths) == 0 {
		return nil
	}

	out := make([]Path, 0, len(paths))
	for _, path := range paths {
		out = append(out, Path(path.String()))
	}

	return out
}
