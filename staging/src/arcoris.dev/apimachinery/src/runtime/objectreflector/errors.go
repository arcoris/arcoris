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

package objectreflector

import (
	"errors"
	"fmt"
)

var (
	// ErrNilSource reports construction with no objectstorewatch source.
	ErrNilSource = errors.New("nil object reflector source")
	// ErrNilSink reports construction with no destination sink.
	ErrNilSink = errors.New("nil object reflector sink")
	// ErrInvalidOption reports a malformed option supplied to New.
	ErrInvalidOption = errors.New("invalid object reflector option")
	// ErrAlreadyRunning reports a concurrent Run call on the same Reflector.
	ErrAlreadyRunning = errors.New("object reflector already running")
	// ErrInvalidEvent reports a malformed watch event observed while reflecting.
	ErrInvalidEvent = errors.New("invalid object reflector event")
	// ErrSourceContractViolation reports invalid ListerWatcher or Source API
	// behavior outside the watch event stream itself.
	ErrSourceContractViolation = errors.New("object reflector source contract violation")
	// ErrChangeOutsideCollection reports a changed event that does not belong to
	// the collection this Reflector owns.
	ErrChangeOutsideCollection = errors.New("object reflector change outside collection")
	// ErrNonMonotonicRevision reports a changed or progress event that moves
	// backward or repeats a revision boundary.
	ErrNonMonotonicRevision = errors.New("object reflector non-monotonic revision")
)

var errRelistRequired = errors.New("object reflector relist required")

// invalidOptionError preserves ErrInvalidOption while keeping the concrete
// reason visible in test output and logs.
func invalidOptionError(cause error) error {
	return errors.Join(ErrInvalidOption, cause)
}

// invalidEventError preserves ErrInvalidEvent for broad source-protocol
// handling while retaining the lower validation or contract cause.
func invalidEventError(cause error) error {
	return errors.Join(ErrInvalidEvent, cause)
}

// changeOutsideCollectionError annotates which committed change escaped the
// requested structural collection.
func changeOutsideCollectionError(cause error) error {
	return errors.Join(ErrChangeOutsideCollection, cause)
}

// nonMonotonicRevisionError annotates the revision boundary that violated the
// stream ordering contract.
func nonMonotonicRevisionError(cause error) error {
	return errors.Join(ErrNonMonotonicRevision, cause)
}

// sourceContractError classifies ListerWatcher/Source behavior outside the
// event stream contract. Event ordering and request matching are delegated to
// objectwatch.Validator instead.
func sourceContractError(format string, args ...any) error {
	return errors.Join(ErrSourceContractViolation, fmt.Errorf(format, args...))
}
