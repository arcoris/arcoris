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

import (
	"context"
	"net/http"

	"arcoris.dev/component-base/pkg/health"
)

// evaluator is the evaluator capability required by Handler.
//
// The public constructor accepts *health.Evaluator because healthhttp is an
// adapter over package health, not a generic health abstraction framework. The
// private interface keeps Handler focused on the single capability it needs and
// makes package-level tests able to exercise evaluator error paths without
// changing the public API.
type evaluator interface {
	Evaluate(ctx context.Context, target health.Target) (health.Report, error)
}

// Handler adapts one health target to http.Handler.
//
// A Handler is bound to exactly one health.Target. It does not own a Registry,
// does not register checks, does not create an Evaluator, does not run periodic
// probes, and does not make restart, routing, admission, scheduling, logging, or
// metrics decisions.
//
// ServeHTTP is safe for concurrent use after construction, provided the
// underlying Evaluator is safe for concurrent use. Package health Evaluator is
// designed to read from Registry at evaluation time and return caller-owned
// reports.
type Handler struct {
	evaluator evaluator
	target    health.Target
	config    config
}

// NewHandler returns an HTTP handler for target.
//
// evaluator MUST be non-nil. target MUST be concrete. Options are applied in
// order after target-specific defaults are selected. Invalid options are returned
// at construction time so ServeHTTP does not have to handle malformed rendering
// configuration.
//
// The handler passes the request context to Evaluator. It does not add its own
// timeout because timeout ownership belongs to health.Evaluator.
func NewHandler(evaluator *health.Evaluator, target health.Target, opts ...Option) (*Handler, error) {
	if evaluator == nil {
		return nil, ErrNilEvaluator
	}
	if !target.IsConcrete() {
		return nil, health.InvalidTargetError{Target: target}
	}

	config := defaultConfig(target)
	if err := applyOptions(&config, opts...); err != nil {
		return nil, err
	}

	return &Handler{
		evaluator: evaluator,
		target:    target,
		config:    config,
	}, nil
}

// ServeHTTP evaluates the configured health target and writes a safe HTTP
// response.
//
// Supported methods are GET and HEAD. Unsupported methods receive 405 Method Not
// Allowed with an Allow header. HEAD receives the same status and headers as GET,
// but render.go suppresses the response body.
//
// Evaluator errors are treated as adapter/evaluator boundary failures and are
// rendered with the configured error status. Raw evaluator errors are not exposed
// in the response body.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h == nil || h.evaluator == nil {
		renderHandlerError(w, r, defaultConfig(health.TargetUnknown))
		return
	}

	if !methodAllowed(r.Method) {
		writeMethodNotAllowed(w)
		return
	}

	report, err := h.evaluator.Evaluate(r.Context(), h.target)
	if err != nil {
		renderHandlerError(w, r, h.config)
		return
	}

	passed := report.Passed(h.config.policy)
	renderReport(w, r, h.config, report, passed)
}

// Target returns the health target served by h.
//
// Target returns health.TargetUnknown for a nil handler. This keeps diagnostic
// code defensive and avoids panics when inspecting optional handler values.
func (h *Handler) Target() health.Target {
	if h == nil {
		return health.TargetUnknown
	}

	return h.target
}

// Format returns the response format configured for h.
//
// Format returns FormatText for a nil handler because text is the safe zero-value
// response format.
func (h *Handler) Format() Format {
	if h == nil {
		return FormatText
	}

	return h.config.format
}

// DetailLevel returns the response detail level configured for h.
//
// DetailLevel returns DetailNone for a nil handler because no details is the safe
// zero-value exposure mode.
func (h *Handler) DetailLevel() DetailLevel {
	if h == nil {
		return DetailNone
	}

	return h.config.detailLevel
}
