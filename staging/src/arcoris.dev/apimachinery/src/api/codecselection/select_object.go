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

// SelectObjectDecoder selects a byte-based object decoder for contentType.
func (s Selector) SelectObjectDecoder(contentType ContentType) (Selection, codec.ObjectDecoder, error) {
	path := runtimeDecodePath(codec.TargetObject.String())
	selection, err := s.selectDecode(contentType, codec.TargetObject, TransportBytes, path)
	if err != nil {
		return Selection{}, nil, err
	}

	selected, ok := selection.Entry.Codec().(codec.ObjectCodec)
	if !ok {
		return Selection{}, nil, missingRuntimeCapabilityError(path, selection)
	}

	return selection, selected, nil
}

// SelectObjectEncoder selects a byte-based object encoder from preferences.
func (s Selector) SelectObjectEncoder(preferences PreferenceSet) (Selection, codec.ObjectEncoder, error) {
	path := runtimeEncodePath(codec.TargetObject.String())
	selection, err := s.selectEncode(preferences, codec.TargetObject, TransportBytes, path)
	if err != nil {
		return Selection{}, nil, err
	}

	selected, ok := selection.Entry.Codec().(codec.ObjectCodec)
	if !ok {
		return Selection{}, nil, missingRuntimeCapabilityError(path, selection)
	}

	return selection, selected, nil
}

// SelectObjectStreamDecoder selects a streaming object decoder for contentType.
func (s Selector) SelectObjectStreamDecoder(contentType ContentType) (Selection, codec.ObjectStreamDecoder, error) {
	path := runtimeDecodePath(codec.TargetObject.String())
	selection, err := s.selectDecode(contentType, codec.TargetObject, TransportStream, path)
	if err != nil {
		return Selection{}, nil, err
	}

	selected, ok := selection.Entry.Codec().(codec.ObjectStreamCodec)
	if !ok {
		return Selection{}, nil, missingRuntimeCapabilityError(path, selection)
	}

	return selection, selected, nil
}

// SelectObjectStreamEncoder selects a streaming object encoder from preferences.
func (s Selector) SelectObjectStreamEncoder(preferences PreferenceSet) (Selection, codec.ObjectStreamEncoder, error) {
	path := runtimeEncodePath(codec.TargetObject.String())
	selection, err := s.selectEncode(preferences, codec.TargetObject, TransportStream, path)
	if err != nil {
		return Selection{}, nil, err
	}

	selected, ok := selection.Entry.Codec().(codec.ObjectStreamCodec)
	if !ok {
		return Selection{}, nil, missingRuntimeCapabilityError(path, selection)
	}

	return selection, selected, nil
}
