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

import "arcoris.dev/health"

// Install creates a health HTTP handler for target and registers it at path.
//
// The function validates the mux, path, evaluator, target, and adapter options
// before mutating mux. Path ownership remains router-neutral: healthhttp checks
// only broad HTTP-path safety invariants and leaves duplicate route behavior to
// the concrete mux implementation.
func Install(mux Mux, path string, evaluator Evaluator, target health.Target, opts ...Option) error {
	if nilMux(mux) {
		return ErrNilMux
	}
	if err := ValidatePath(path); err != nil {
		return err
	}

	handler, err := NewHandler(evaluator, target, opts...)
	if err != nil {
		return err
	}

	mux.Handle(path, handler)
	return nil
}

// InstallDefaults registers the default startup, liveness, and readiness health
// endpoints on mux.
//
// Compatibility paths such as "/healthz" and "/health" are not installed by
// default because they do not have universal target semantics.
//
// The function constructs every handler before registering any route. If the
// evaluator or any option is invalid, mux is left untouched and no partial
// installation occurs.
//
// Options apply to every installed endpoint. Callers that need different
// per-target policy or representation settings should call Install directly for
// each route instead of relying on InstallDefaults.
func InstallDefaults(mux Mux, evaluator Evaluator, opts ...Option) error {
	if nilMux(mux) {
		return ErrNilMux
	}

	handlers := make([]defaultHandler, 0, len(defaultHandlers))
	for _, item := range defaultHandlers {
		handler, err := NewHandler(evaluator, item.target, opts...)
		if err != nil {
			return err
		}

		handlers = append(handlers, defaultHandler{path: item.path, target: item.target, handler: handler})
	}

	for _, item := range handlers {
		mux.Handle(item.path, item.handler)
	}

	return nil
}

// defaultHandler describes one default health HTTP endpoint registration.
//
// The slice below is treated as immutable package configuration.
type defaultHandler struct {
	path    string
	target  health.Target
	handler *Handler
}

// defaultHandlers lists the endpoints installed by InstallDefaults.
var defaultHandlers = []defaultHandler{
	{path: DefaultStartupPath, target: health.TargetStartup},
	{path: DefaultLivePath, target: health.TargetLive},
	{path: DefaultReadyPath, target: health.TargetReady},
}
