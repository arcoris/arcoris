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

// SelectValueDecoder selects a byte-based value decoder for contentType.
func (s Selector) SelectValueDecoder(contentType ContentType) (Selection, codec.ValueDecoder, error) {
	path := runtimeDecodePath(codec.TargetValue.String())
	selection, err := s.selectDecode(contentType, codec.TargetValue, TransportBytes, path)
	if err != nil {
		return Selection{}, nil, err
	}

	selected, ok := selection.Entry.Codec().(codec.ValueCodec)
	if !ok {
		return Selection{}, nil, missingRuntimeCapabilityError(path, selection)
	}

	return selection, selected, nil
}

// SelectValueEncoder selects a byte-based value encoder from preferences.
func (s Selector) SelectValueEncoder(preferences PreferenceSet) (Selection, codec.ValueEncoder, error) {
	path := runtimeEncodePath(codec.TargetValue.String())
	selection, err := s.selectEncode(preferences, codec.TargetValue, TransportBytes, path)
	if err != nil {
		return Selection{}, nil, err
	}

	selected, ok := selection.Entry.Codec().(codec.ValueCodec)
	if !ok {
		return Selection{}, nil, missingRuntimeCapabilityError(path, selection)
	}

	return selection, selected, nil
}

// SelectValueStreamDecoder selects a streaming value decoder for contentType.
func (s Selector) SelectValueStreamDecoder(contentType ContentType) (Selection, codec.ValueStreamDecoder, error) {
	path := runtimeDecodePath(codec.TargetValue.String())
	selection, err := s.selectDecode(contentType, codec.TargetValue, TransportStream, path)
	if err != nil {
		return Selection{}, nil, err
	}

	selected, ok := selection.Entry.Codec().(codec.ValueStreamCodec)
	if !ok {
		return Selection{}, nil, missingRuntimeCapabilityError(path, selection)
	}

	return selection, selected, nil
}

// SelectValueStreamEncoder selects a streaming value encoder from preferences.
func (s Selector) SelectValueStreamEncoder(preferences PreferenceSet) (Selection, codec.ValueStreamEncoder, error) {
	path := runtimeEncodePath(codec.TargetValue.String())
	selection, err := s.selectEncode(preferences, codec.TargetValue, TransportStream, path)
	if err != nil {
		return Selection{}, nil, err
	}

	selected, ok := selection.Entry.Codec().(codec.ValueStreamCodec)
	if !ok {
		return Selection{}, nil, missingRuntimeCapabilityError(path, selection)
	}

	return selection, selected, nil
}
