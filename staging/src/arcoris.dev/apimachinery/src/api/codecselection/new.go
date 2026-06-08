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

// New validates config and returns an immutable selector.
//
// Construction is all-or-nothing. Bindings are normalized, checked for
// duplicate keys, and validated against the supplied codecregistry.Registry.
func New(config Config) (Selector, error) {
	decode, decodeBindings, err := buildDecodeBindings(config)
	if err != nil {
		return Selector{}, err
	}

	encode, encodeBindings, err := buildEncodeBindings(config)
	if err != nil {
		return Selector{}, err
	}

	return Selector{
		registry:       config.Registry,
		decode:         decode,
		encode:         encode,
		decodeBindings: decodeBindings,
		encodeBindings: encodeBindings,
	}, nil
}
