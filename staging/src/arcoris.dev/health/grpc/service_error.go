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

package healthgrpc

import (
	"errors"
	"fmt"
)

var (
	// ErrInvalidService identifies an invalid gRPC health service mapping name.
	ErrInvalidService = errors.New("healthgrpc: invalid service")

	// ErrDuplicateService identifies a repeated gRPC health service mapping name.
	ErrDuplicateService = errors.New("healthgrpc: duplicate service")
)

// InvalidServiceError describes an invalid service name in a health mapping.
//
// Service names are configuration identifiers, so the value may appear in the
// error. Reason is a stable package-owned diagnostic string, not a raw lower
// level error.
type InvalidServiceError struct {
	// Service is the rejected gRPC service name.
	Service string

	// Index is the mapping index where validation failed.
	Index int

	// Reason is the package-owned validation reason.
	Reason string
}

// Error returns a stable diagnostic string for an invalid service mapping.
func (e InvalidServiceError) Error() string {
	return fmt.Sprintf("%v: service=%q index=%d reason=%q", ErrInvalidService, e.Service, e.Index, e.Reason)
}

// Is reports compatibility with ErrInvalidService.
func (e InvalidServiceError) Is(target error) bool {
	return target == ErrInvalidService
}

// DuplicateServiceError describes a repeated service name in the configured
// service mapping list.
type DuplicateServiceError struct {
	// Service is the repeated gRPC service name.
	Service string

	// Index is the duplicate mapping index.
	Index int

	// PreviousIndex is the earlier mapping index with the same service name.
	PreviousIndex int
}

// Error returns a stable diagnostic string for a duplicate service mapping.
func (e DuplicateServiceError) Error() string {
	return fmt.Sprintf(
		"%v: service=%q index=%d previous_index=%d",
		ErrDuplicateService,
		e.Service,
		e.Index,
		e.PreviousIndex,
	)
}

// Is reports compatibility with ErrDuplicateService.
func (e DuplicateServiceError) Is(target error) bool {
	return target == ErrDuplicateService
}
