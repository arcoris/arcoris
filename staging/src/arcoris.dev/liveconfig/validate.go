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

// Validator checks a normalized candidate configuration before publication.
//
// A validator MUST NOT mutate its argument. Use a Normalizer for defaulting or
// canonicalization. When validation returns an error, Holder keeps the previous
// last-good configuration active and does not advance the source revision.
//
// Validators should check semantic invariants of the final candidate: bounds,
// required fields, mutually exclusive fields, and relationships between values.
// They should not read external sources or perform side effects, because Apply
// holds the holder write mutex while validation runs.
type Validator[T any] func(T) error

// validateValue applies cfg's validator when one is configured.
//
// A nil validator means all normalized candidates are accepted by this package.
// Domain code may still enforce stricter guarantees before calling Apply.
func validateValue[T any](cfg config[T], val T) error {
	if cfg.validate == nil {
		return nil
	}

	return cfg.validate(val)
}
