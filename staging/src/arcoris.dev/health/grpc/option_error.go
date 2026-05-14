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
	"time"
)

var (
	// ErrNilSource reports a nil Source passed to NewServer.
	ErrNilSource = errors.New("healthgrpc: nil source")

	// ErrNilRegistrar reports a nil grpc.ServiceRegistrar passed to Register.
	ErrNilRegistrar = errors.New("healthgrpc: nil registrar")

	// ErrNilServer reports a nil Server passed to Register.
	ErrNilServer = errors.New("healthgrpc: nil server")

	// ErrNilOption reports a nil Option passed to NewServer.
	ErrNilOption = errors.New("healthgrpc: nil option")

	// ErrNilClock reports a nil clock passed through WithClock or config.
	ErrNilClock = errors.New("healthgrpc: nil clock")

	// ErrInvalidWatchInterval classifies non-positive Watch intervals.
	ErrInvalidWatchInterval = errors.New("healthgrpc: invalid watch interval")

	// ErrInvalidMaxListServices classifies non-positive List service limits.
	ErrInvalidMaxListServices = errors.New("healthgrpc: invalid max list services")
)

// InvalidWatchIntervalError describes an invalid Watch polling interval.
type InvalidWatchIntervalError struct {
	// Interval is the rejected interval value.
	Interval time.Duration
}

// Error returns a stable diagnostic string for an invalid Watch interval.
func (e InvalidWatchIntervalError) Error() string {
	return fmt.Sprintf("%v: interval=%s", ErrInvalidWatchInterval, e.Interval)
}

// Is reports compatibility with ErrInvalidWatchInterval.
func (e InvalidWatchIntervalError) Is(target error) bool {
	return target == ErrInvalidWatchInterval
}

// InvalidMaxListServicesError describes an invalid List service-count limit.
type InvalidMaxListServicesError struct {
	// Max is the rejected maximum service count.
	Max int
}

// Error returns a stable diagnostic string for an invalid List service limit.
func (e InvalidMaxListServicesError) Error() string {
	return fmt.Sprintf("%v: max=%d", ErrInvalidMaxListServices, e.Max)
}

// Is reports compatibility with ErrInvalidMaxListServices.
func (e InvalidMaxListServicesError) Is(target error) bool {
	return target == ErrInvalidMaxListServices
}

// validateWatchInterval rejects intervals that cannot drive a polling ticker.
func validateWatchInterval(interval time.Duration) error {
	if interval <= 0 {
		return InvalidWatchIntervalError{Interval: interval}
	}

	return nil
}

// validateMaxListServices rejects limits that would make List unusable.
func validateMaxListServices(max int) error {
	if max <= 0 {
		return InvalidMaxListServicesError{Max: max}
	}

	return nil
}
