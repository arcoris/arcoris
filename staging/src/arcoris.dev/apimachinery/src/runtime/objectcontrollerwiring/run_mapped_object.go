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

import "context"

// RunMappedObject coordinates one already assembled mapped-object controller
// graph.
//
// RunMappedObject has the same lifecycle policy as RunSameObject: it starts the
// reflector and controller together, shuts down the queue when either side
// exits, cancels the sibling component, waits for both sides, and preserves
// fatal component errors.
//
// The function panics when ctx is nil, matching the lower-level runtime
// components. It validates only the handles it coordinates directly; the
// assembled controller already owns its snapshot source dependency.
func RunMappedObject(ctx context.Context, graph *MappedObject) error {
	if ctx == nil {
		panic("objectcontrollerwiring: nil context")
	}
	if err := validateRunnableMappedObject(graph); err != nil {
		return err
	}

	return runGraph(ctx, graph.Queue(), graph.Reflector(), graph.Controller())
}

// validateRunnableMappedObject keeps runner validation narrow. A graph is
// runnable when the queue can be shut down, the reflector can be canceled, and
// the controller can drain the queue.
func validateRunnableMappedObject(graph *MappedObject) error {
	if graph == nil ||
		graph.Queue() == nil ||
		graph.Reflector() == nil ||
		graph.Controller() == nil {
		return ErrInvalidMappedObject
	}

	return nil
}
