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

// selectTypedDecoder shares exact binding lookup and runtime capability checks.
func selectTypedDecoder[T any](
	s Selector,
	contentType ContentType,
	target codec.Target,
	transport Transport,
) (Selection, T, error) {
	path := runtimeDecodePath(target.String())
	selection, err := s.selectDecode(contentType, target, transport, path)
	if err != nil {
		var zero T

		return Selection{}, zero, err
	}

	return typedSelection[T](selection, path)
}

// selectTypedEncoder shares preference lookup and runtime capability checks.
func selectTypedEncoder[T any](
	s Selector,
	preferences PreferenceSet,
	target codec.Target,
	transport Transport,
) (Selection, T, error) {
	path := runtimeEncodePath(target.String())
	selection, err := s.selectEncode(preferences, target, transport, path)
	if err != nil {
		var zero T

		return Selection{}, zero, err
	}

	return typedSelection[T](selection, path)
}

// typedSelection defensively re-asserts the selected runtime capability.
func typedSelection[T any](selection Selection, path string) (Selection, T, error) {
	selected, ok := any(selection.Entry.Codec()).(T)
	if !ok {
		var zero T

		return Selection{}, zero, missingRuntimeCapabilityError(path, selection)
	}

	return selection, selected, nil
}
