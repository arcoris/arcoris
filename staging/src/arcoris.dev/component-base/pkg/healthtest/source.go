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

package healthtest

import (
	"context"
	"errors"
	"sync"

	"arcoris.dev/component-base/pkg/health"
)

// SourceFunc adapts a function to an evaluator-shaped health report source.
//
// The shape intentionally matches health.Evaluator.Evaluate and adapter source
// interfaces such as healthgrpc.Source without importing any adapter package.
type SourceFunc func(context.Context, health.Target) (health.Report, error)

// Evaluate calls fn with ctx and target.
//
// SourceFunc performs no nil guard. A nil function should fail where the test
// invokes it, which keeps the failure close to the fixture setup mistake.
func (fn SourceFunc) Evaluate(ctx context.Context, target health.Target) (health.Report, error) {
	return fn(ctx, target)
}

// StaticSource returns the same report for every target.
//
// It is useful when a transport adapter test only cares about status mapping or
// error handling and does not need target-specific source behavior.
type StaticSource struct {
	// report is stored with a detached Checks slice.
	report health.Report
}

// NewStaticSource returns a Source that always returns report.
//
// The report is copied at construction and on every Evaluate call. Tests may
// mutate their original report or the returned report without affecting later
// source evaluations.
func NewStaticSource(report health.Report) *StaticSource {
	return &StaticSource{report: copyReport(report)}
}

// Evaluate returns the configured report.
//
// target is ignored by design. Use TargetSource or CountingSource when the test
// needs target-specific behavior.
func (s *StaticSource) Evaluate(context.Context, health.Target) (health.Report, error) {
	if s == nil {
		return health.Report{}, nil
	}

	return copyReport(s.report), nil
}

// TargetSource returns reports by target.
//
// Missing targets return UnknownReport(target) instead of an error. This makes
// it convenient for adapter tests that treat absence of source data as an
// unknown health observation rather than a transport failure.
type TargetSource struct {
	// mu protects reports for concurrent streaming/watch tests.
	mu sync.RWMutex

	// reports stores detached reports by target.
	reports map[health.Target]health.Report
}

// NewTargetSource returns a Source backed by reports.
//
// reports is copied deeply enough for health.Report.Checks so caller mutations
// do not affect later evaluations.
func NewTargetSource(reports map[health.Target]health.Report) *TargetSource {
	return &TargetSource{reports: copyReportMap(reports)}
}

// Evaluate returns the report configured for target or an unknown report.
//
// The returned report is detached from the source's internal map.
func (s *TargetSource) Evaluate(_ context.Context, target health.Target) (health.Report, error) {
	if s == nil {
		return UnknownReport(target), nil
	}

	s.mu.RLock()
	report, ok := s.reports[target]
	s.mu.RUnlock()
	if !ok {
		return UnknownReport(target), nil
	}

	return copyReport(report), nil
}

// ErrorSource always returns an evaluation error.
//
// It is intended for adapter tests that verify generic public error handling
// and non-leakage of raw source errors.
type ErrorSource struct {
	// err is returned by every Evaluate call.
	err error
}

// NewErrorSource returns a Source that always fails with err.
//
// A nil err is replaced with a stable package-owned error so the source remains
// useful even when the test does not care about the exact raw cause.
func NewErrorSource(err error) *ErrorSource {
	if err == nil {
		err = errors.New("healthtest source error")
	}

	return &ErrorSource{err: err}
}

// Evaluate returns the configured error.
//
// The returned report is always zero because Source failures should be handled
// by the adapter path under test, not by inspecting partial reports.
func (s *ErrorSource) Evaluate(context.Context, health.Target) (health.Report, error) {
	if s == nil || s.err == nil {
		return health.Report{}, errors.New("healthtest source error")
	}

	return health.Report{}, s.err
}

// CountingSource records per-target calls while returning configured reports.
//
// CountingSource is concurrency-safe and useful for List/watch tests that need
// to verify deduplication or polling behavior without relying on sleeps.
type CountingSource struct {
	// mu protects reports, errors, and calls.
	mu sync.Mutex

	// reports stores the report returned for each target.
	reports map[health.Target]health.Report

	// errors stores target-specific evaluation failures.
	errors map[health.Target]error

	// calls stores per-target evaluation counts.
	calls map[health.Target]int
}

// NewCountingSource returns a call-counting Source backed by reports.
//
// Reports are copied at construction and on read. Callers may mutate their input
// map or returned reports without changing source state.
func NewCountingSource(reports map[health.Target]health.Report) *CountingSource {
	return &CountingSource{
		reports: copyReportMap(reports),
		errors:  make(map[health.Target]error),
		calls:   make(map[health.Target]int),
	}
}

// SetReport configures the report returned for target.
//
// SetReport is safe to call before or during tests that evaluate the source, but
// deterministic tests should still prefer configuring all reports up front.
func (s *CountingSource) SetReport(target health.Target, report health.Report) {
	if s == nil {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.reports[target] = copyReport(report)
}

// SetError configures an evaluation error for target.
//
// Passing nil removes the target-specific error and restores report/unknown
// behavior for that target.
func (s *CountingSource) SetError(target health.Target, err error) {
	if s == nil {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if err == nil {
		delete(s.errors, target)
		return
	}
	s.errors[target] = err
}

// Evaluate records target and returns its configured report or error.
//
// Missing targets return UnknownReport(target). That default keeps tests focused
// on source behavior without forcing every fixture to enumerate all concrete
// targets.
func (s *CountingSource) Evaluate(_ context.Context, target health.Target) (health.Report, error) {
	if s == nil {
		return UnknownReport(target), nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.calls[target]++
	if err := s.errors[target]; err != nil {
		return health.Report{}, err
	}
	if report, ok := s.reports[target]; ok {
		return copyReport(report), nil
	}

	return UnknownReport(target), nil
}

// Calls returns how many times target has been evaluated.
//
// Calls is safe to use while another goroutine is evaluating the source.
func (s *CountingSource) Calls(target health.Target) int {
	if s == nil {
		return 0
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	return s.calls[target]
}

// SequenceSource returns a sticky sequence of reports for one target.
//
// SequenceSource is useful for watch/polling tests. It advances only when the
// requested target matches its configured target, and the final report remains
// sticky after the sequence is exhausted.
type SequenceSource struct {
	// target is the only target that consumes the report sequence.
	target health.Target

	// mu protects reports and calls across streaming tests.
	mu sync.Mutex

	// reports is the remaining sticky report sequence.
	reports []health.Report

	// calls records evaluations by requested target.
	calls map[health.Target]int
}

// NewSequenceSource returns a Source backed by a report sequence for target.
//
// The report slice and each report's Checks slice are copied so caller
// mutations cannot rewrite the sequence.
func NewSequenceSource(target health.Target, reports ...health.Report) *SequenceSource {
	copied := make([]health.Report, len(reports))
	for i, report := range reports {
		copied[i] = copyReport(report)
	}

	return &SequenceSource{
		target:  target,
		reports: copied,
		calls:   make(map[health.Target]int),
	}
}

// Evaluate returns the next report for the configured target.
//
// UnknownReport(target) is returned for other targets and for empty sequences.
// That behavior keeps target mistakes observable without manufacturing raw
// source errors.
func (s *SequenceSource) Evaluate(_ context.Context, target health.Target) (health.Report, error) {
	if s == nil {
		return UnknownReport(target), nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.calls[target]++
	if target != s.target || len(s.reports) == 0 {
		return UnknownReport(target), nil
	}

	report := s.reports[0]
	if len(s.reports) > 1 {
		s.reports = s.reports[1:]
	}

	return copyReport(report), nil
}

// Calls returns how many times target has been evaluated.
//
// Calls includes evaluations for non-sequenced targets, which helps tests assert
// that adapters asked for the intended target.
func (s *SequenceSource) Calls(target health.Target) int {
	if s == nil {
		return 0
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	return s.calls[target]
}

// copyReport returns report with a detached Checks slice.
//
// health.Report is a plain value with a mutable Checks slice. Sources copy that
// slice at every boundary so tests cannot accidentally share mutable fixture
// state across evaluations.
func copyReport(report health.Report) health.Report {
	report.Checks = report.ChecksCopy()
	return report
}

// copyReportMap returns reports with detached map and report slices.
//
// The map is copied even when empty so source constructors can safely store and
// later mutate their internal maps without aliasing caller input.
func copyReportMap(reports map[health.Target]health.Report) map[health.Target]health.Report {
	copied := make(map[health.Target]health.Report, len(reports))
	for target, report := range reports {
		copied[target] = copyReport(report)
	}

	return copied
}
