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
	return selectTypedDecoder[codec.ValueCodec](s, contentType, codec.TargetValue, TransportBytes)
}

// SelectValueEncoder selects a byte-based value encoder from preferences.
func (s Selector) SelectValueEncoder(preferences PreferenceSet) (Selection, codec.ValueEncoder, error) {
	return selectTypedEncoder[codec.ValueCodec](s, preferences, codec.TargetValue, TransportBytes)
}

// SelectValueStreamDecoder selects a streaming value decoder for contentType.
func (s Selector) SelectValueStreamDecoder(contentType ContentType) (Selection, codec.ValueStreamDecoder, error) {
	return selectTypedDecoder[codec.ValueStreamCodec](s, contentType, codec.TargetValue, TransportStream)
}

// SelectValueStreamEncoder selects a streaming value encoder from preferences.
func (s Selector) SelectValueStreamEncoder(preferences PreferenceSet) (Selection, codec.ValueStreamEncoder, error) {
	return selectTypedEncoder[codec.ValueStreamCodec](s, preferences, codec.TargetValue, TransportStream)
}
