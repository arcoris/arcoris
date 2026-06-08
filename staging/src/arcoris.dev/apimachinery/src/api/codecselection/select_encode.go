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

import "arcoris.dev/apimachinery/api/codec"

// selectEncode returns the first supported encode binding from preferences.
func (s Selector) selectEncode(
	preferences PreferenceSet,
	target codec.Target,
	transport Transport,
	path string,
) (Selection, error) {
	normalized, err := normalizePreferenceSetAt(path+".preferences", preferences.items)
	if err != nil {
		return Selection{}, err
	}
	if normalized.IsZero() {
		return Selection{}, errorAt(
			path+".preferences",
			ErrNoEncodePreference,
			ErrorReasonNoEncodePreference,
			"encode preferences are required",
		)
	}

	for _, preference := range normalized.items {
		record, ok := s.encode[bindingKey{
			contentType: preference.contentType.key(),
			target:      target,
			transport:   transport,
		}]
		if ok {
			return selectionFromRecord(DirectionEncode, record), nil
		}
	}

	return Selection{}, errorfAt(
		path,
		ErrNoEncodePreference,
		ErrorReasonNoEncodePreference,
		"no encode binding matched %d preferences for target %q and transport %q",
		normalized.Len(),
		target,
		transport,
	)
}
