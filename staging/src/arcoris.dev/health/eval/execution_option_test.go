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

package eval

import (
	"errors"
	"testing"

	"arcoris.dev/health"
)

func TestEvaluatorConfigDefaultExecutionPolicy(t *testing.T) {
	t.Parallel()

	config := defaultEvaluatorConfig()
	want := DefaultExecutionPolicy()

	if config.executionPolicy != want {
		t.Fatalf("execution policy = %+v, want %+v", config.executionPolicy, want)
	}
	if config.targetExecutionPolicies == nil {
		t.Fatal("targetExecutionPolicies is nil")
	}
}

func TestWithExecutionPolicy(t *testing.T) {
	t.Parallel()

	config := defaultEvaluatorConfig()
	err := WithExecutionPolicy(ParallelExecutionPolicy(4))(&config)
	if err != nil {
		t.Fatalf("WithExecutionPolicy() = %v, want nil", err)
	}

	want := ParallelExecutionPolicy(4)
	if config.executionPolicy != want {
		t.Fatalf("execution policy = %+v, want %+v", config.executionPolicy, want)
	}
}

func TestWithExecutionPolicyRejectsInvalidPolicy(t *testing.T) {
	t.Parallel()

	config := defaultEvaluatorConfig()
	err := WithExecutionPolicy(ParallelExecutionPolicy(0))(&config)

	if !errors.Is(err, ErrInvalidExecutionPolicy) {
		t.Fatalf("WithExecutionPolicy(invalid) = %v, want ErrInvalidExecutionPolicy", err)
	}
}

func TestWithTargetExecutionPolicy(t *testing.T) {
	t.Parallel()

	config := defaultEvaluatorConfig()
	err := WithTargetExecutionPolicy(health.TargetReady, ParallelExecutionPolicy(4))(&config)
	if err != nil {
		t.Fatalf("WithTargetExecutionPolicy() = %v, want nil", err)
	}

	want := ParallelExecutionPolicy(4)
	if got := config.targetExecutionPolicies[health.TargetReady]; got != want {
		t.Fatalf("target execution policy = %+v, want %+v", got, want)
	}
}

func TestWithTargetExecutionPolicyRejectsInvalidTarget(t *testing.T) {
	t.Parallel()

	config := defaultEvaluatorConfig()
	err := WithTargetExecutionPolicy(health.TargetUnknown, ParallelExecutionPolicy(4))(&config)

	if !errors.Is(err, health.ErrInvalidTarget) {
		t.Fatalf("WithTargetExecutionPolicy(invalid target) = %v, want health.ErrInvalidTarget", err)
	}
}

func TestWithTargetExecutionPolicyRejectsInvalidPolicy(t *testing.T) {
	t.Parallel()

	config := defaultEvaluatorConfig()
	err := WithTargetExecutionPolicy(health.TargetReady, ParallelExecutionPolicy(0))(&config)

	if !errors.Is(err, ErrInvalidExecutionPolicy) {
		t.Fatalf("WithTargetExecutionPolicy(invalid policy) = %v, want ErrInvalidExecutionPolicy", err)
	}
}

func TestExecutionOptionConveniences(t *testing.T) {
	t.Parallel()

	config := defaultEvaluatorConfig()
	err := applyEvaluatorOptions(
		&config,
		WithParallelChecks(8),
		WithTargetParallelChecks(health.TargetReady, 4),
		WithTargetSequentialChecks(health.TargetLive),
	)
	if err != nil {
		t.Fatalf("applyEvaluatorOptions() = %v, want nil", err)
	}

	if got, want := config.executionPolicy, ParallelExecutionPolicy(8); got != want {
		t.Fatalf("default execution policy = %+v, want %+v", got, want)
	}
	if got, want := config.targetExecutionPolicies[health.TargetReady], ParallelExecutionPolicy(4); got != want {
		t.Fatalf("ready execution policy = %+v, want %+v", got, want)
	}
	if got, want := config.targetExecutionPolicies[health.TargetLive], DefaultExecutionPolicy(); got != want {
		t.Fatalf("live execution policy = %+v, want %+v", got, want)
	}
}

func TestExecutionOptionsApplyInOrder(t *testing.T) {
	t.Parallel()

	config := defaultEvaluatorConfig()
	err := applyEvaluatorOptions(
		&config,
		WithParallelChecks(2),
		WithParallelChecks(4),
		WithTargetParallelChecks(health.TargetReady, 2),
		WithTargetParallelChecks(health.TargetReady, 6),
	)
	if err != nil {
		t.Fatalf("applyEvaluatorOptions() = %v, want nil", err)
	}

	if got, want := config.executionPolicy, ParallelExecutionPolicy(4); got != want {
		t.Fatalf("default execution policy = %+v, want %+v", got, want)
	}
	if got, want := config.targetExecutionPolicies[health.TargetReady], ParallelExecutionPolicy(6); got != want {
		t.Fatalf("ready execution policy = %+v, want %+v", got, want)
	}
}
