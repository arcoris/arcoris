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
	return selectTypedDecoder[codec.ObjectCodec](s, contentType, codec.TargetObject, TransportBytes)
}

// SelectObjectEncoder selects a byte-based object encoder from preferences.
func (s Selector) SelectObjectEncoder(preferences PreferenceSet) (Selection, codec.ObjectEncoder, error) {
	return selectTypedEncoder[codec.ObjectCodec](s, preferences, codec.TargetObject, TransportBytes)
}

// SelectObjectStreamDecoder selects a streaming object decoder for contentType.
func (s Selector) SelectObjectStreamDecoder(contentType ContentType) (Selection, codec.ObjectStreamDecoder, error) {
	return selectTypedDecoder[codec.ObjectStreamCodec](s, contentType, codec.TargetObject, TransportStream)
}

// SelectObjectStreamEncoder selects a streaming object encoder from preferences.
func (s Selector) SelectObjectStreamEncoder(preferences PreferenceSet) (Selection, codec.ObjectStreamEncoder, error) {
	return selectTypedEncoder[codec.ObjectStreamCodec](s, preferences, codec.TargetObject, TransportStream)
}
