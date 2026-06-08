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

package bulkhead

import "errors"

var (
	// ErrNilBulkhead reports a method call on a nil Bulkhead receiver.
	ErrNilBulkhead = errors.New("bulkhead: nil bulkhead")

	// ErrUninitializedBulkhead reports use of a zero Bulkhead value instead of a
	// value created by New.
	ErrUninitializedBulkhead = errors.New("bulkhead: uninitialized bulkhead")

	// ErrInvalidLease reports a method call on a nil or zero Lease instead of a
	// lease returned by a successful acquisition.
	ErrInvalidLease = errors.New("bulkhead: invalid lease")

	// ErrLeaseReleased reports a strict Release call after the lease has already
	// returned its capacity.
	ErrLeaseReleased = errors.New("bulkhead: lease already released")
)
