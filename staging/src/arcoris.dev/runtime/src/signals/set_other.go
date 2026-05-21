//go:build !unix && !windows

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

// ShutdownSignals returns the portable interrupt signal for platforms without
// Unix or Windows-specific process signal sets.
//
// The result is deliberately conservative. It exposes only os.Interrupt because
// this low-level package should not invent platform policy when the standard
// library does not define a conventional shutdown set.
func ShutdownSignals() []os.Signal {
	return []os.Signal{os.Interrupt}
}

// ReloadSignals returns no default reload signals on this platform.
//
// Reload behavior is opt-in. Callers that have a platform-specific reload
// signal should pass it explicitly.
func ReloadSignals() []os.Signal {
	return nil
}

// DiagnosticSignals returns no default diagnostic signals on this platform.
//
// Diagnostic behavior is owner policy. This package exposes platform
// conventions only when the platform provides one.
func DiagnosticSignals() []os.Signal {
	return nil
}
