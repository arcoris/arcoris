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

package snapshot

import "arcoris.dev/chrono/clock"

// passiveClock returns the Publisher clock.
//
// A zero-value Publisher has no configured clock, so it lazily falls back to
// clock.RealClock. NewPublisher should be used when deterministic timestamps are
// required in tests.
func (p *Publisher[T]) passiveClock() clock.PassiveClock {
	if p.clock != nil {
		return p.clock
	}

	return clock.RealClock{}
}
