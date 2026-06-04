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

package codec

// Supports reports whether i declares support for target.
//
// Supports is a simple metadata query and intentionally does not validate the
// Info first. Invalid or duplicate target declarations are handled by Validate
// and Normalize.
func (i Info) Supports(target Target) bool {
	for _, supported := range i.Targets {
		if supported == target {
			return true
		}
	}

	return false
}

// SupportsMediaType reports whether i declares support for mediaType.
//
// The comparison is exact. Callers that accept non-canonical input should
// normalize the media type before querying.
func (i Info) SupportsMediaType(mediaType MediaType) bool {
	for _, supported := range i.MediaTypes {
		if supported == mediaType {
			return true
		}
	}

	return false
}
