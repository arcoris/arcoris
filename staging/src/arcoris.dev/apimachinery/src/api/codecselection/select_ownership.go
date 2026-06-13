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
	return selectTypedDecoder[codec.ObjectOwnershipCodec](s, contentType, codec.TargetObjectOwnership, TransportBytes)
}

// SelectObjectOwnershipEncoder selects a byte-based ownership encoder from preferences.
func (s Selector) SelectObjectOwnershipEncoder(
	preferences PreferenceSet,
) (Selection, codec.ObjectOwnershipEncoder, error) {
	return selectTypedEncoder[codec.ObjectOwnershipCodec](s, preferences, codec.TargetObjectOwnership, TransportBytes)
}

// SelectObjectOwnershipStreamDecoder selects a streaming ownership decoder for contentType.
func (s Selector) SelectObjectOwnershipStreamDecoder(
	contentType ContentType,
) (Selection, codec.ObjectOwnershipStreamDecoder, error) {
	return selectTypedDecoder[codec.ObjectOwnershipStreamCodec](s, contentType, codec.TargetObjectOwnership, TransportStream)
}

// SelectObjectOwnershipStreamEncoder selects a streaming ownership encoder from preferences.
func (s Selector) SelectObjectOwnershipStreamEncoder(
	preferences PreferenceSet,
) (Selection, codec.ObjectOwnershipStreamEncoder, error) {
	return selectTypedEncoder[codec.ObjectOwnershipStreamCodec](s, preferences, codec.TargetObjectOwnership, TransportStream)
}
