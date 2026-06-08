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

// normalizeContentTypeAt validates c and returns its normalized form.
func normalizeContentTypeAt(path string, c ContentType) (ContentType, error) {
	mediaType, err := c.mediaType.Normalize()
	if err != nil {
		return ContentType{}, wrapAt(
			path+".mediaType",
			ErrInvalidContentType,
			ErrorReasonInvalidContentType,
			"content type media type is invalid",
			err,
		)
	}

	parameters, err := normalizeParametersAt(path+".parameters", c.parameters.items)
	if err != nil {
		return ContentType{}, err
	}

	return ContentType{mediaType: mediaType, parameters: parameters}, nil
}
