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

// Preference is one normalized encode content type preference.
type Preference struct {
	// contentType is the normalized requested content type.
	contentType ContentType

	// weight ranks this preference against siblings.
	weight Weight

	// order preserves caller order for deterministic equal-weight tie-breaking.
	order int
}

// NewPreference validates contentType and weight into one encode preference.
func NewPreference(contentType ContentType, weight Weight) (Preference, error) {
	return normalizePreferenceAt(pathPreference, Preference{contentType: contentType, weight: weight})
}

// MustPreference returns a normalized Preference or panics when input is invalid.
func MustPreference(contentType ContentType, weight Weight) Preference {
	preference, err := NewPreference(contentType, weight)
	if err != nil {
		panic(err)
	}

	return preference
}

// IsZero reports whether p contains no content type and no weight.
func (p Preference) IsZero() bool {
	return p.contentType.IsZero() && p.weight.IsZero()
}

// ContentType returns the normalized preferred content type.
func (p Preference) ContentType() ContentType {
	return p.contentType
}

// Weight returns the validated preference weight.
func (p Preference) Weight() Weight {
	return p.weight
}
