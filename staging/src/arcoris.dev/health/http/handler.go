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
	"net/http"

	"arcoris.dev/health"
)

// Handler adapts one concrete health target to an HTTP endpoint.
//
// Handler is intentionally thin. It does not own health registry contents,
// check execution policy, retries, logging, metrics, tracing, authentication,
// authorization, routing, or restart decisions. Its responsibility is limited
// to invoking a health.Evaluator for one target and rendering the resulting
// report through the package's safe HTTP surface.
//
// Handler stores only adapter configuration. It never exposes Result.Cause,
// panic stacks, raw errors, or context causes to callers.
type Handler struct {
	evaluator *health.Evaluator
	target    health.Target
	config    config
}

// NewHandler returns an HTTP handler for target.
//
// evaluator MUST be non-nil. target MUST be concrete. Options configure only
// HTTP-adapter policy such as response format, safe detail level, and status
// code mapping.
func NewHandler(evaluator *health.Evaluator, target health.Target, opts ...Option) (*Handler, error) {
	if evaluator == nil {
		return nil, ErrNilEvaluator
	}
	if !target.IsConcrete() {
		return nil, health.InvalidTargetError{Target: target}
	}

	cfg := defaultConfig(target)
	if err := applyOptions(&cfg, opts...); err != nil {
		return nil, err
	}

	return &Handler{evaluator: evaluator, target: target, config: cfg}, nil
}

// ServeHTTP evaluates the configured health target and writes a safe response.
//
// The method accepts only GET and HEAD. Unsupported methods receive a generic
// 405 response with the same cache and content-safety headers used by normal
// health responses.
//
// A nil handler, nil evaluator, or nil request is treated as an adapter-boundary
// failure and returns a generic handler error response instead of panicking.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h == nil || h.evaluator == nil || r == nil {
		renderHandlerError(w, r, defaultConfig(health.TargetUnknown))
		return
	}

	if !methodAllowed(r.Method) {
		writeMethodNotAllowed(w, r)
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

// Target returns the concrete health target served by h.
//
// A nil handler reports TargetUnknown so diagnostics and tests can inspect the
// zero boundary without panicking.
func (h *Handler) Target() health.Target {
	if h == nil {
		return health.TargetUnknown
	}

	return h.target
}

// Format returns the response format configured for h.
//
// A nil handler reports FormatText because text is the package's zero/default
// representation.
func (h *Handler) Format() Format {
	if h == nil {
		return FormatText
	}

	return h.config.format
}

// DetailLevel returns the response detail level configured for h.
//
// A nil handler reports DetailNone because safe-by-default responses do not
// include per-check details unless the owner opts in explicitly.
func (h *Handler) DetailLevel() DetailLevel {
	if h == nil {
		return DetailNone
	}

	return h.config.detailLevel
}
