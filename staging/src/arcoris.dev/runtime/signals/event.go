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

import "os"

// Event describes one process signal observed by this package.
//
// Event is a copyable value object. It contains only the received os.Signal and
// intentionally does not include timestamps, counters, process identifiers,
// escalation flags, lifecycle state, logging fields, or metrics labels. Those
// details belong to the runtime owner that interprets the event.
type Event struct {
	// Signal is the process signal that was received.
	//
	// Signal is non-nil for Events produced by this package. The zero Event is a
	// valid empty value but does not describe an observed signal.
	Signal os.Signal
}
