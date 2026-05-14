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

package probe

import "errors"

var (
	// ErrNilEvaluator identifies a nil health evaluator passed to NewRunner.
	//
	// Runner calls Evaluator.Evaluate during every probe cycle. A nil evaluator
	// would make the runner unable to produce any health observation.
	ErrNilEvaluator = errors.New("healthprobe: nil evaluator")

	// ErrNilRunner identifies a Run call on a nil Runner receiver.
	//
	// Snapshot and Snapshots treat nil Runner receivers as empty readers. Run
	// returns this stable error because a nil Runner cannot own a probe loop.
	ErrNilRunner = errors.New("healthprobe: nil runner")

	// ErrRunnerRunning identifies a concurrent Run call on the same Runner.
	//
	// Runner owns one schedule-driven loop. Running two loops for the same Runner
	// would duplicate probe work and create ambiguous snapshot generations.
	ErrRunnerRunning = errors.New("healthprobe: runner already running")
)
