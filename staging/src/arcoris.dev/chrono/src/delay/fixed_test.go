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

package delay

import (
	"testing"
	"time"
)

func TestFixedRejectsNegativeDelay(t *testing.T) {
	mustPanicWith(t, errNegativeFixedDelay, func() {
		Fixed(-time.Nanosecond)
	})
}

func TestFixedReturnsConfiguredDelayForever(t *testing.T) {
	seq := Fixed(3 * time.Second).NewSequence()

	mustNext(t, seq, 3*time.Second)
	mustNext(t, seq, 3*time.Second)
	mustNext(t, seq, 3*time.Second)
}

func TestFixedAllowsZeroDelay(t *testing.T) {
	seq := Fixed(0).NewSequence()

	mustNext(t, seq, 0)
	mustNext(t, seq, 0)
}
