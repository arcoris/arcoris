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

	"arcoris.dev/apimachinery/api/objectwatch"
)

// processEvent validates one source event before routing it by kind.
//
// The clone at the boundary keeps a misbehaving stream implementation from
// sharing mutable event payload state with the sink. Only EventChanged reaches
// Sink; progress and restart-required events are reflector control flow.
func (r *Reflector) processEvent(ctx context.Context, event objectwatch.Event) error {
	event = event.Clone()
	if err := event.Validate(); err != nil {
		return invalidEventError(err)
	}

	switch event.Kind {
	case objectwatch.EventChanged:
		return r.processChanged(ctx, event)
	case objectwatch.EventProgress:
		return r.processProgress(event)
	case objectwatch.EventRestartRequired:
		return relistRequiredError(event.Restart)
	default:
		return sourceContractError("unknown event kind %s", event.Kind.String())
	}
}
