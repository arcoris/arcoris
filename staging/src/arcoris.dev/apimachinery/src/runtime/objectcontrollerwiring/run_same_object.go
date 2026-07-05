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

import (
	"context"
)

// RunSameObject coordinates one already assembled same-object controller graph.
//
// The runner starts the reflector and controller with a shared child context.
// When either side exits, it closes the producer/consumer boundary by shutting
// down the queue, canceling the sibling side, and waiting for both components to
// return. The queue shutdown is intentionally graph-level policy: the
// controller knows how to drain a shut-down queue, but it does not decide when
// the producer side is finished.
//
// Fatal component errors are preserved. Parent cancellation is returned only
// when no more specific component failure exists.
//
// RunSameObject panics when ctx is nil. It does not recover panics from lower
// components.
func RunSameObject(ctx context.Context, graph *SameObject) error {
	if ctx == nil {
		panic("objectcontrollerwiring: nil context")
	}
	if err := validateRunnableSameObject(graph); err != nil {
		return err
	}

	return runGraph(ctx, graph.Queue(), graph.Reflector(), graph.Controller())
}

// validateRunnableSameObject checks the components whose lifecycle is directly
// coordinated by RunSameObject. Cache is intentionally not checked here: the
// assembled controller already owns its snapshot source dependency, while the
// runner only needs the queue, reflector, and controller handles.
func validateRunnableSameObject(graph *SameObject) error {
	if graph == nil ||
		graph.Queue() == nil ||
		graph.Reflector() == nil ||
		graph.Controller() == nil {
		return ErrInvalidSameObject
	}

	return nil
}
