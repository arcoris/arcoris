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

package codec

// EncodeOptions controls common format-independent encode behavior.
//
// Zero EncodeOptions are valid and request the implementation's normal output.
// These options describe byte formatting only; they do not request value
// normalization, metadata mutation, field ownership updates, or validation.
type EncodeOptions struct {
	// Deterministic requests stable output when the format supports it.
	//
	// Deterministic output should preserve semantic content while making byte
	// order stable for tests, signatures, caching, or reproducible artifacts.
	Deterministic bool

	// Pretty requests human-readable output without changing semantic content.
	//
	// Pretty formatting is advisory. Binary formats may ignore it or return an
	// unsupported-feature error if they choose to reject the request explicitly.
	Pretty bool
}
