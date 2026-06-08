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

// selectDecode returns the validated decode binding for exact key material.
func (s Selector) selectDecode(
	contentType ContentType,
	target codec.Target,
	transport Transport,
	path string,
) (Selection, error) {
	normalized, err := normalizeContentTypeAt(path+".contentType", contentType)
	if err != nil {
		return Selection{}, wrapAt(
			path+".contentType",
			ErrInvalidContentType,
			ErrorReasonInvalidContentType,
			"decode content type is invalid",
			err,
		)
	}

	record, ok := s.decode[bindingKey{
		contentType: normalized.key(),
		target:      target,
		transport:   transport,
	}]
	if !ok {
		return Selection{}, errorfAt(
			path,
			ErrNoDecodeBinding,
			ErrorReasonNoDecodeBinding,
			"no decode binding for content type %q, target %q, and transport %q",
			normalized,
			target,
			transport,
		)
	}

	return selectionFromRecord(DirectionDecode, record), nil
}
