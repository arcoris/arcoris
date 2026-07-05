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

package objectcontrollerwiring

import "errors"

var (
	// ErrInvalidSameObject reports a SameObject graph that is missing a component
	// required by RunSameObject to coordinate reflector, queue, and controller
	// lifecycle.
	ErrInvalidSameObject = errors.New("objectcontrollerwiring: invalid same-object graph")

	// ErrInvalidMappedObject reports a MappedObject graph that is missing a
	// component required by RunMappedObject to coordinate reflector, queue, and
	// controller lifecycle.
	ErrInvalidMappedObject = errors.New("objectcontrollerwiring: invalid mapped-object graph")

	// ErrInvalidMultiSource reports a MultiSource graph that is missing a
	// component required by RunMultiSource to coordinate all reflectors, the
	// shared queue, and the controller lifecycle.
	ErrInvalidMultiSource = errors.New("objectcontrollerwiring: invalid multi-source graph")

	// ErrNilIndex reports a nil prebuilt object index in graph input config.
	// Index construction remains owned by runtime/objectindex; wiring only
	// validates that configured index sinks can be installed into fanout order.
	ErrNilIndex = errors.New("objectcontrollerwiring: nil index")
)
