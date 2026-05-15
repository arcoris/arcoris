/*
  Copyright 2026 The ARCORIS Authors

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
*/

package liveconfig

// Normalizer converts a candidate configuration into its canonical form.
//
// A normalizer may apply defaults, canonicalize names, normalize nil and empty
// collections, derive fields, or reject a candidate that cannot be normalized.
// It runs after cloning and before validation.
//
// Normalization is part of the candidate transaction. If it returns an error,
// Holder keeps the previous last-good value, records the error as LastError for
// Apply, and does not advance the revision.
type Normalizer[T any] func(T) (T, error)

// normalizeValue applies cfg's normalizer when one is configured.
//
// A nil normalizer means the cloned candidate is already canonical enough for
// the package to validate and publish.
func normalizeValue[T any](cfg config[T], val T) (T, error) {
	if cfg.normalize == nil {
		return val, nil
	}

	return cfg.normalize(val)
}
