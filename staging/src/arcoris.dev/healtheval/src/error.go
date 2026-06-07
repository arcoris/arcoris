// Copyright 2026 The ARCORIS Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package eval

import (
	"errors"
	"fmt"

	"arcoris.dev/health"
)

var (
	// ErrNilResolver identifies a nil health.CheckResolver passed to NewEvaluator.
	//
	// Evaluator requires a resolver owner. A nil resolver would make target
	// evaluation ambiguous, so construction rejects it instead of producing an
	// evaluator that fails later.
	ErrNilResolver = errors.New("healtheval: nil check resolver")

	// ErrNilEvaluatorOption identifies a nil option passed to NewEvaluator.
	//
	// Options are executed during evaluator construction. A nil option is a
	// wiring error and is rejected explicitly.
	ErrNilEvaluatorOption = errors.New("healtheval: nil evaluator option")

	// ErrNilClock identifies a nil clock passed to WithClock.
	//
	// Evaluator uses clock.PassiveClock for observation timestamps and elapsed
	// duration measurement. A nil clock would panic during evaluation.
	ErrNilClock = errors.New("healtheval: nil clock")

	// ErrInvalidTimeout identifies a negative health check timeout.
	//
	// A zero timeout is valid and disables evaluator-enforced timeout. Negative
	// values do not describe a meaningful evaluation budget.
	ErrInvalidTimeout = errors.New("healtheval: invalid timeout")

	// ErrInvalidCheckResult identifies a structurally invalid health.Result
	// returned by a checker during evaluator-owned normalization.
	//
	// Evaluator converts invalid checker output into an unknown misconfiguration
	// result instead of returning an invalid health.Report to callers.
	ErrInvalidCheckResult = errors.New("healtheval: invalid check result")

	// ErrMismatchedCheckResult identifies a health.Result whose non-empty name
	// differs from the health.Checker name that produced it.
	//
	// health.Checker identity is owned by health.Checker.Name and registry
	// registration. Evaluator rejects mismatched result names so one checker
	// cannot accidentally or intentionally publish another checker's identity in
	// reports.
	ErrMismatchedCheckResult = errors.New("healtheval: mismatched check result")

	// ErrMismatchedResolvedTarget identifies a resolver that returned a CheckSet
	// for a different target than the target being evaluated.
	//
	// This is an evaluator-owned resolver contract violation. Evaluator returns
	// an unknown report for the requested target and does not execute checks from
	// the mismatched set.
	ErrMismatchedResolvedTarget = errors.New("healtheval: mismatched resolved target")
)

// InvalidCheckResultError describes an invalid health.Result returned by a
// checker.
//
// InvalidCheckResultError is classified as ErrInvalidCheckResult. Callers should
// use errors.Is for classification and inspect fields only for diagnostics.
type InvalidCheckResultError struct {
	CheckName string
	Result    health.Result
}

// Error returns the invalid check result message.
func (e InvalidCheckResultError) Error() string {
	return fmt.Sprintf(
		"%v: check=%q status=%s reason=%s duration=%s",
		ErrInvalidCheckResult,
		e.CheckName,
		e.Result.Status.String(),
		e.Result.Reason.String(),
		e.Result.Duration,
	)
}

// Is reports whether target matches the invalid check result classification.
func (e InvalidCheckResultError) Is(target error) bool {
	return target == ErrInvalidCheckResult
}

// MismatchedCheckResultError describes a checker result whose name does not
// match the checker that produced it.
//
// MismatchedCheckResultError is classified as ErrMismatchedCheckResult. Callers
// should use errors.Is for classification and inspect fields only for
// diagnostics.
type MismatchedCheckResultError struct {
	CheckName  string
	ResultName string
}

// Error returns the mismatched check result message.
func (e MismatchedCheckResultError) Error() string {
	return fmt.Sprintf(
		"%v: check=%q result=%q",
		ErrMismatchedCheckResult,
		e.CheckName,
		e.ResultName,
	)
}

// Is reports whether target matches the mismatched check result classification.
func (e MismatchedCheckResultError) Is(target error) bool {
	return target == ErrMismatchedCheckResult
}

// MismatchedResolvedTargetError describes a resolver result that belongs to a
// different target than the target requested by Evaluator.
//
// MismatchedResolvedTargetError is classified as ErrMismatchedResolvedTarget.
// Callers should use errors.Is for classification and inspect fields only for
// diagnostics.
type MismatchedResolvedTargetError struct {
	// Requested is the concrete target passed to Evaluator.Evaluate.
	Requested health.Target

	// Resolved is the target carried by the returned health.CheckSet.
	Resolved health.Target
}

// Error returns the mismatched resolved target message.
func (e MismatchedResolvedTargetError) Error() string {
	return fmt.Sprintf(
		"%v: requested=%s resolved=%s",
		ErrMismatchedResolvedTarget,
		e.Requested.String(),
		e.Resolved.String(),
	)
}

// Is reports whether target matches the mismatched resolved target
// classification.
func (e MismatchedResolvedTargetError) Is(target error) bool {
	return target == ErrMismatchedResolvedTarget
}
