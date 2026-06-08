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

// SelectObjectOwnershipDecoder selects a byte-based ownership decoder for contentType.
func (s Selector) SelectObjectOwnershipDecoder(
	contentType ContentType,
) (Selection, codec.ObjectOwnershipDecoder, error) {
	path := runtimeDecodePath(codec.TargetObjectOwnership.String())
	selection, err := s.selectDecode(contentType, codec.TargetObjectOwnership, TransportBytes, path)
	if err != nil {
		return Selection{}, nil, err
	}

	selected, ok := selection.Entry.Codec().(codec.ObjectOwnershipCodec)
	if !ok {
		return Selection{}, nil, missingRuntimeCapabilityError(path, selection)
	}

	return selection, selected, nil
}

// SelectObjectOwnershipEncoder selects a byte-based ownership encoder from preferences.
func (s Selector) SelectObjectOwnershipEncoder(
	preferences PreferenceSet,
) (Selection, codec.ObjectOwnershipEncoder, error) {
	path := runtimeEncodePath(codec.TargetObjectOwnership.String())
	selection, err := s.selectEncode(preferences, codec.TargetObjectOwnership, TransportBytes, path)
	if err != nil {
		return Selection{}, nil, err
	}

	selected, ok := selection.Entry.Codec().(codec.ObjectOwnershipCodec)
	if !ok {
		return Selection{}, nil, missingRuntimeCapabilityError(path, selection)
	}

	return selection, selected, nil
}

// SelectObjectOwnershipStreamDecoder selects a streaming ownership decoder for contentType.
func (s Selector) SelectObjectOwnershipStreamDecoder(
	contentType ContentType,
) (Selection, codec.ObjectOwnershipStreamDecoder, error) {
	path := runtimeDecodePath(codec.TargetObjectOwnership.String())
	selection, err := s.selectDecode(contentType, codec.TargetObjectOwnership, TransportStream, path)
	if err != nil {
		return Selection{}, nil, err
	}

	selected, ok := selection.Entry.Codec().(codec.ObjectOwnershipStreamCodec)
	if !ok {
		return Selection{}, nil, missingRuntimeCapabilityError(path, selection)
	}

	return selection, selected, nil
}

// SelectObjectOwnershipStreamEncoder selects a streaming ownership encoder from preferences.
func (s Selector) SelectObjectOwnershipStreamEncoder(
	preferences PreferenceSet,
) (Selection, codec.ObjectOwnershipStreamEncoder, error) {
	path := runtimeEncodePath(codec.TargetObjectOwnership.String())
	selection, err := s.selectEncode(preferences, codec.TargetObjectOwnership, TransportStream, path)
	if err != nil {
		return Selection{}, nil, err
	}

	selected, ok := selection.Entry.Codec().(codec.ObjectOwnershipStreamCodec)
	if !ok {
		return Selection{}, nil, missingRuntimeCapabilityError(path, selection)
	}

	return selection, selected, nil
}
