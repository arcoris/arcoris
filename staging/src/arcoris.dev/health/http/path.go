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

package healthhttp

const (
	// DefaultStartupPath is the default HTTP path for startup health evaluation.
	//
	// Startup answers whether a component has completed enough bootstrapping to
	// be considered started. This path maps to health.TargetStartup in default
	// installs.
	DefaultStartupPath = "/startupz"

	// DefaultLivePath is the default HTTP path for liveness health evaluation.
	//
	// Liveness answers whether a component is alive and able to continue making
	// progress. This path maps to health.TargetLive in default installs.
	DefaultLivePath = "/livez"

	// DefaultReadyPath is the default HTTP path for readiness health evaluation.
	//
	// Readiness answers whether a component should receive new work. This path
	// maps to health.TargetReady in default installs.
	DefaultReadyPath = "/readyz"

	// DefaultHealthPath is the canonical HTTP path for general health evaluation.
	//
	// General health answers whether the component is broadly healthy according
	// to the caller's bound health target or aggregate handler. It is
	// intentionally separate from startup, liveness, and readiness paths, which
	// have narrower operational meanings.
	//
	// The package provides this path as the canonical general health route, but
	// callers still decide which handler or target is bound to it.
	DefaultHealthPath = "/healthz"
)
