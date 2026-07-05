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

	"arcoris.dev/apimachinery/runtime/objectreflector"
)

// RunMultiSource coordinates one already assembled multi-source controller
// graph.
//
// RunMultiSource starts every input reflector and the single controller under a
// shared child context. When any component exits, it shuts down the shared
// queue, cancels the remaining components, waits for all of them, and returns
// the fatal component error set or the parent context error.
func RunMultiSource(ctx context.Context, graph *MultiSource) error {
	if ctx == nil {
		panic("objectcontrollerwiring: nil context")
	}
	if err := validateRunnableMultiSource(graph); err != nil {
		return err
	}

	return runComponents(ctx, graph.Queue(), graph.Controller(), multiSourceReflectors(graph)...)
}

// validateRunnableMultiSource checks the handles directly coordinated by the
// runner. Input caches are not checked here because the assembled controller
// already owns the primary snapshot source dependency, and secondary caches are
// not part of lifecycle coordination.
func validateRunnableMultiSource(graph *MultiSource) error {
	if graph == nil ||
		graph.Queue() == nil ||
		graph.Controller() == nil ||
		graph.Primary() == nil ||
		graph.Primary().Reflector() == nil {
		return ErrInvalidMultiSource
	}
	for _, input := range graph.Secondary() {
		if input == nil || input.Reflector() == nil {
			return ErrInvalidMultiSource
		}
	}

	return nil
}

// multiSourceReflectors returns the reflectors in input order. That order is
// used only to start goroutines deterministically; all reflectors feed the same
// shared queue.
func multiSourceReflectors(graph *MultiSource) []*objectreflector.Reflector {
	inputs := graph.Inputs()
	reflectors := make([]*objectreflector.Reflector, 0, len(inputs))
	for _, input := range inputs {
		reflectors = append(reflectors, input.Reflector())
	}

	return reflectors
}
