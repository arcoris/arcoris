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

// normalizePreferenceAt validates p at path.
func normalizePreferenceAt(path string, p Preference) (Preference, error) {
	contentType, err := normalizeContentTypeAt(path+".contentType", p.contentType)
	if err != nil {
		return Preference{}, wrapAt(
			path+".contentType",
			ErrInvalidPreference,
			ErrorReasonInvalidPreference,
			"preference content type is invalid",
			err,
		)
	}
	if err := p.weight.Validate(); err != nil {
		return Preference{}, wrapAt(
			path+".weight",
			ErrInvalidPreference,
			ErrorReasonInvalidPreference,
			"preference weight is invalid",
			err,
		)
	}

	return Preference{contentType: contentType, weight: p.weight, order: p.order}, nil
}
