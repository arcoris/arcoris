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

package objectworkqueue

import "errors"

var (
	// ErrInvalidCapacity reports that Options.Capacity is not positive.
	ErrInvalidCapacity = errors.New("objectworkqueue: invalid capacity")

	// ErrInvalidQueue reports a nil Queue receiver on an error-returning method.
	ErrInvalidQueue = errors.New("objectworkqueue: invalid queue")

	// ErrInvalidItem reports an item whose object key is structurally invalid.
	ErrInvalidItem = errors.New("objectworkqueue: invalid item")

	// ErrFull reports that TryAdd cannot add a new distinct item because the
	// queue is at capacity.
	ErrFull = errors.New("objectworkqueue: full")

	// ErrShutDown reports that the queue is shutting down or has already shut
	// down.
	ErrShutDown = errors.New("objectworkqueue: shut down")

	// ErrNotProcessing reports Done for an item that is not currently
	// processing.
	ErrNotProcessing = errors.New("objectworkqueue: item is not processing")
)
