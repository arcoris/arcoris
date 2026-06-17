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

package objectwatch

// Validate checks the watch start invariant for its mode.
func (s Start) Validate() error {
	if s.IsZero() {
		return invalidStartError("zero start")
	}
	if !s.Mode.IsValid() {
		return invalidStartError("mode %s is invalid", s.Mode.String())
	}

	switch s.Mode {
	case StartAfterRevision:
		return nil
	case StartAtCurrent:
		if !s.Revision.IsZero() {
			return invalidStartError("atCurrent requires zero revision")
		}
	}

	return nil
}
