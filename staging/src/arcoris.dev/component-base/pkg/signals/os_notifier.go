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

package signals

import (
	"os"
	"os/signal"
)

// osNotifier adapts package os/signal to the package-local notifier seam.
type osNotifier struct{}

// notify registers ch to receive the supplied process signals.
func (osNotifier) notify(ch chan<- os.Signal, sigs ...os.Signal) {
	signal.Notify(ch, sigs...)
}

// stop unregisters ch from package os/signal delivery.
func (osNotifier) stop(ch chan<- os.Signal) {
	signal.Stop(ch)
}
