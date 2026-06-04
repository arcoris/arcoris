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

import "io"

// targetDecoder converts a parsed JSON document root into one codec target.
type targetDecoder[T any] func(jsonPath, jsonNode, resolvedDecodeConfig) (T, error)

// decodeTargetFrom owns the common stream-to-node-to-target decode pipeline.
//
// The public DecodeValue/DecodeObject/DecodeObjectOwnership methods stay
// explicit, while this helper keeps strict parsing, UTF-8 checks, duplicate-key
// checks, trailing-data checks, and depth handling identical for every target.
func decodeTargetFrom[T any](
	r io.Reader,
	config resolvedDecodeConfig,
	decode targetDecoder[T],
) (T, error) {
	node, err := decodeJSONDocument(r, config)
	if err != nil {
		var zero T
		return zero, err
	}

	return decode(rootPath(), node, config)
}
