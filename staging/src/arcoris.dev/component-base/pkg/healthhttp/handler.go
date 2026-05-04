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
type evaluator interface {
	Evaluate(ctx context.Context, target health.Target) (health.Report, error)
}

// Handler adapts one health target to http.Handler.
type Handler struct {
	evaluator evaluator
	target    health.Target
	config    config
}

// NewHandler returns an HTTP handler for target.
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

	return &Handler{evaluator: evaluator, target: target, config: config}, nil
}

// ServeHTTP evaluates the configured health target and writes a safe response.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h == nil || h.evaluator == nil || r == nil {
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
func (h *Handler) Target() health.Target {
	if h == nil {
		return health.TargetUnknown
	}

	return h.target
}

// Format returns the response format configured for h.
func (h *Handler) Format() Format {
	if h == nil {
		return FormatText
	}

	return h.config.format
}

// DetailLevel returns the response detail level configured for h.
func (h *Handler) DetailLevel() DetailLevel {
	if h == nil {
		return DetailNone
	}

	return h.config.detailLevel
}
