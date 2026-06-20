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
	"context"
	"errors"

	"arcoris.dev/apimachinery/api/objectwatch"
)

// consumeStream drains one watch stream until context cancellation, sink
// failure, source terminal error, or a relist-required event ends the cycle.
//
// The reflector deliberately keeps no internal event buffer here. If a sink is
// slow, the source stream sees that backpressure directly and can report
// continuity loss according to the objectwatch contract.
func (r *Reflector) consumeStream(
	ctx context.Context,
	stream objectwatch.Stream,
	validator *objectwatch.Validator,
) error {
	for {
		event, err := stream.Next(ctx)
		if err != nil {
			return err
		}
		event = event.Clone()
		if err := validator.Accept(event); err != nil {
			if errors.Is(err, objectwatch.ErrInvalidEvent) {
				return invalidEventError(err)
			}
			return err
		}
		if err := r.processEvent(ctx, event); err != nil {
			return err
		}
	}
}
