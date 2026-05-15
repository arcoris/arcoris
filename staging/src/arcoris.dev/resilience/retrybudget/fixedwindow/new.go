/*
  Copyright 2026 The ARCORIS Authors

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
*/

package fixedwindow

import (
	"arcoris.dev/resilience/retrybudget"
	"arcoris.dev/snapshot"
)

// New creates a fixed-window retry budget limiter.
//
// New validates opts, starts the first accounting window at the configured
// clock's current time, publishes an initial snapshot, and returns the limiter.
func New(opts ...Option) (*Limiter, error) {
	cfg, err := newConfig(opts...)
	if err != nil {
		return nil, err
	}

	publisher := snapshot.NewPublisher[retrybudget.Snapshot](snapshot.WithClock(cfg.clock))

	l := &Limiter{
		cfg:         cfg,
		windowStart: cfg.clock.Now(),
		published:   publisher,
	}

	l.mu.Lock()
	l.publishLocked()
	l.mu.Unlock()

	return l, nil
}
