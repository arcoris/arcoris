//go:build unix

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
	"syscall"
)

// ShutdownSignals returns the default graceful-shutdown signal set for Unix
// systems.
//
// The set contains os.Interrupt and SIGTERM, which are the process-level
// triggers ordinary component entrypoints should treat as owner-requested
// shutdown. Each call returns a fresh slice so callers may append, sort, or
// otherwise normalize the result without mutating package-owned state.
func ShutdownSignals() []os.Signal {
	return []os.Signal{os.Interrupt, syscall.SIGTERM}
}

// ReloadSignals returns the conventional configuration-reload signal set for
// Unix systems.
//
// Reload support is opt-in. The package never registers reload signals through
// NotifyContext or ShutdownController defaults because reload policy belongs to
// the component owner.
func ReloadSignals() []os.Signal {
	return []os.Signal{syscall.SIGHUP}
}

// DiagnosticSignals returns the conventional diagnostic signal set for Unix
// systems.
//
// Diagnostic signals are opt-in. This package only exposes the platform
// convention; it does not decide whether diagnostics are dumped, logged,
// exported, or ignored.
func DiagnosticSignals() []os.Signal {
	return []os.Signal{syscall.SIGQUIT}
}
