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

package healthprobe

import (
	"errors"
	"fmt"

	"arcoris.dev/component-base/pkg/health"
)

var (
	// ErrNoTargets identifies a missing probe target configuration.
	ErrNoTargets = errors.New("healthprobe: no targets")

	// ErrDuplicateTarget identifies a repeated target in one probe target list.
	ErrDuplicateTarget = errors.New("healthprobe: duplicate target")
)

// DuplicateTargetError describes a repeated health target in a probe target list.
type DuplicateTargetError struct {
	Target        health.Target
	Index         int
	PreviousIndex int
}

func (e DuplicateTargetError) Error() string {
	return fmt.Sprintf(
		"%v: target=%s index=%d previous_index=%d",
		ErrDuplicateTarget,
		e.Target.String(),
		e.Index,
		e.PreviousIndex,
	)
}

func (e DuplicateTargetError) Is(target error) bool {
	return target == ErrDuplicateTarget
}
