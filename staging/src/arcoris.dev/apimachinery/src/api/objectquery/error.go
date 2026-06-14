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

package objectquery

import "errors"

var (
	// ErrInvalidQuery classifies malformed top-level query values.
	ErrInvalidQuery = errors.New("invalid object query")

	// ErrInvalidSelector classifies malformed metadata selector values.
	ErrInvalidSelector = errors.New("invalid object query selector")

	// ErrInvalidRequirement classifies malformed metadata requirements.
	ErrInvalidRequirement = errors.New("invalid object query requirement")

	// ErrInvalidOperator classifies unknown requirement operators.
	ErrInvalidOperator = errors.New("invalid object query operator")
)
