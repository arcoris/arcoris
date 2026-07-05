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

package objectindex

import "errors"

var (
	// ErrInvalidIndex reports a nil or internally malformed Index, or an invalid
	// index operation such as an empty lookup name or emitted value.
	ErrInvalidIndex = errors.New("objectindex: invalid index")

	// ErrInvalidDefinition reports malformed index definitions passed to New.
	ErrInvalidDefinition = errors.New("objectindex: invalid definition")

	// ErrUnknownIndex reports lookup against a name that was not registered when
	// the Index was constructed.
	ErrUnknownIndex = errors.New("objectindex: unknown index")

	// ErrNilExtractor reports a nil extractor in a definition or nil
	// ExtractorFunc invocation.
	ErrNilExtractor = errors.New("objectindex: nil extractor")
)

// errorWith preserves the local error category while keeping the lower-level
// cause visible to errors.Is callers.
func errorWith(category error, cause error) error {
	if cause == nil {
		return category
	}

	return errors.Join(category, cause)
}
