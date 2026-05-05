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

package health

import "testing"

func TestExecutionPolicyNormalize(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		policy ExecutionPolicy
		want   ExecutionPolicy
	}{
		{
			name:   "zero",
			policy: ExecutionPolicy{},
			want: ExecutionPolicy{
				Mode:           ExecutionSequential,
				MaxConcurrency: 1,
			},
		},
		{
			name: "sequential ignores limit",
			policy: ExecutionPolicy{
				Mode:           ExecutionSequential,
				MaxConcurrency: 64,
			},
			want: ExecutionPolicy{
				Mode:           ExecutionSequential,
				MaxConcurrency: 1,
			},
		},
		{
			name: "parallel preserves limit",
			policy: ExecutionPolicy{
				Mode:           ExecutionParallel,
				MaxConcurrency: 4,
			},
			want: ExecutionPolicy{
				Mode:           ExecutionParallel,
				MaxConcurrency: 4,
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			if got := test.policy.Normalize(); got != test.want {
				t.Fatalf("Normalize() = %+v, want %+v", got, test.want)
			}
		})
	}
}

func TestExecutionPolicyValidate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		policy ExecutionPolicy
		want   bool
	}{
		{name: "zero", policy: ExecutionPolicy{}, want: true},
		{name: "default", policy: DefaultExecutionPolicy(), want: true},
		{name: "parallel", policy: ParallelExecutionPolicy(4), want: true},
		{name: "parallel zero limit", policy: ParallelExecutionPolicy(0), want: false},
		{name: "parallel negative limit", policy: ParallelExecutionPolicy(-1), want: false},
		{name: "invalid mode", policy: ExecutionPolicy{Mode: ExecutionMode(99), MaxConcurrency: 1}, want: false},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			err := test.policy.Validate()
			if got := err == nil; got != test.want {
				t.Fatalf("Validate() ok = %v, want %v; err=%v", got, test.want, err)
			}
		})
	}
}
