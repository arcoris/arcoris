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

package stamp

import (
	"encoding/json"
	"strconv"
)

// MarshalText validates and encodes the generation as decimal text.
func (g Generation) MarshalText() ([]byte, error) {
	return marshalText(g.String(), g.Validate)
}

// UnmarshalText decodes and validates a decimal generation.
func (g *Generation) UnmarshalText(data []byte) error {
	if g == nil {
		return nilReceiver("generation")
	}

	value, err := strconv.ParseUint(string(data), 10, 64)
	if err != nil {
		return invalid(
			"generation",
			ErrInvalidGeneration,
			ErrorReasonInvalidForm,
			"expected unsigned decimal generation",
		)
	}

	*g = Generation(value)
	return nil
}

// MarshalJSON validates and encodes the generation as one JSON number.
func (g Generation) MarshalJSON() ([]byte, error) {
	if err := g.Validate(); err != nil {
		return nil, err
	}

	return json.Marshal(uint64(g))
}

// UnmarshalJSON decodes and validates a JSON number generation.
func (g *Generation) UnmarshalJSON(data []byte) error {
	if g == nil {
		return nilReceiver("generation")
	}

	var value uint64
	if err := json.Unmarshal(data, &value); err != nil {
		return &Error{
			Path:   "generation",
			Err:    ErrInvalidJSON,
			Reason: ErrorReasonInvalidJSON,
			Detail: "expected JSON unsigned integer",
			Cause:  err,
		}
	}

	*g = Generation(value)
	return nil
}
