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

import "arcoris.dev/apimachinery/api/codecregistry"

// lookupEntryAt returns the registry entry identified by entryID.
func lookupEntryAt(path string, config Config, entryID codecregistry.EntryID) (codecregistry.Entry, error) {
	entry, ok := config.Registry.LookupID(entryID)
	if !ok {
		return codecregistry.Entry{}, errorfAt(
			path,
			ErrUnknownEntryID,
			ErrorReasonUnknownEntryID,
			"registry entry ID %q is not registered",
			entryID,
		)
	}

	return entry, nil
}
