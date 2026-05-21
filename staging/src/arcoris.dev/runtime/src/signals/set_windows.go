//go:build windows

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

// ShutdownSignals returns the default graceful-shutdown signal set for Windows.
//
// Windows exposes os.Interrupt as the portable owner-requested shutdown signal
// available to this low-level package. Each call returns a fresh slice so
// callers cannot mutate package-owned defaults.
func ShutdownSignals() []os.Signal {
	return []os.Signal{os.Interrupt}
}

// ReloadSignals returns no default reload signals on Windows.
//
// Reload behavior is intentionally opt-in and platform-specific. Callers that
// support a Windows-specific reload trigger should pass that signal explicitly.
func ReloadSignals() []os.Signal {
	return nil
}

// DiagnosticSignals returns no default diagnostic signals on Windows.
//
// Diagnostic behavior is owner policy. The package does not invent a synthetic
// diagnostic signal when the platform has no conventional default.
func DiagnosticSignals() []os.Signal {
	return nil
}
