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

package codecjson

import (
	"arcoris.dev/apimachinery/api/codec"
	"arcoris.dev/apimachinery/api/codecjson/jsonconfig"
)

var (
	// Compile-time checks keep the concrete codec aligned with api/codec's full
	// byte and streaming aggregate contracts.
	_ codec.Codec          = Codec{}
	_ codec.StreamingCodec = Codec{}
)

// Codec implements JSON byte-slice and stream codecs for all v1 targets.
//
// Codec is immutable after construction and carries only normalized local
// configuration. It does not register itself globally.
type Codec struct {
	// decode stores resolved parser and target-decoder policy.
	decode resolvedDecodeConfig

	// encode stores resolved JSON writer policy.
	encode resolvedEncodeConfig
}

// New constructs a JSON codec.
//
// Public configuration lives in api/codecjson/jsonconfig. New resolves and
// validates that configuration once, then stores private immutable policy so
// runtime encode and decode methods stay optionless.
func New(config jsonconfig.Config) (Codec, error) {
	resolved, err := resolveConfig(config)
	if err != nil {
		return Codec{}, err
	}

	return Codec{
		decode: resolved.decode,
		encode: resolved.encode,
	}, nil
}
