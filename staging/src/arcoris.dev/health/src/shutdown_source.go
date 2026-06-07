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

package health

import "context"

// newSourceChannelCheck constructs a channel-backed shutdown-style checker.
func newSourceChannelCheck(name string, done <-chan struct{}, reason Reason, message string) (Checker, error) {
	if err := ValidateCheckName(name); err != nil {
		return nil, err
	}
	if done == nil {
		return nil, ErrNilSourceChannel
	}

	return sourceChannelCheck{
		name:    name,
		done:    done,
		reason:  reason,
		message: message,
	}, nil
}

// newSourceContextCheck constructs a context-backed shutdown-style checker.
func newSourceContextCheck(name string, source context.Context, reason Reason, message string) (Checker, error) {
	if err := ValidateCheckName(name); err != nil {
		return nil, err
	}
	if source == nil {
		return nil, ErrNilSourceContext
	}

	return sourceContextCheck{
		name:    name,
		source:  source,
		reason:  reason,
		message: message,
	}, nil
}

// sourceChannelCheck reports unhealthy after a channel is closed.
//
// The type is intentionally private. Public callers should construct it through
// NewShutdownCheck or NewDrainCheck so name and channel invariants are validated.
type sourceChannelCheck struct {
	// name is the stable check name returned by Name and attached to results.
	name string

	// done is the owner-published shutdown or drain signal.
	done <-chan struct{}

	// reason classifies the unhealthy result after done is closed.
	reason Reason

	// message is the safe diagnostic text attached after done is closed.
	message string
}

// Name returns the stable check name.
func (c sourceChannelCheck) Name() string {
	return c.name
}

// Check returns healthy while the source channel is open and unhealthy after it
// is closed.
//
// Check ignores ctx because the check is a non-blocking read of owner-published
// state. Evaluator-owned cancellation and timeout are handled by Evaluator.
func (c sourceChannelCheck) Check(ctx context.Context) Result {
	select {
	case <-c.done:
		return Unhealthy(c.name, c.reason, c.message)
	default:
		return Healthy(c.name)
	}
}

// sourceContextCheck reports unhealthy after a source context is cancelled.
//
// The type is intentionally private. Public callers should construct it through
// NewContextShutdownCheck or NewContextDrainCheck so name and context invariants
// are validated.
type sourceContextCheck struct {
	// name is the stable check name returned by Name and attached to results.
	name string

	// source is the owner-published shutdown or drain context observed by Check.
	source context.Context

	// reason classifies the unhealthy result after source is cancelled.
	reason Reason

	// message is the safe diagnostic text attached after source is cancelled.
	message string
}

// Name returns the stable check name.
func (c sourceContextCheck) Name() string {
	return c.name
}

// Check returns healthy while the source context is active and unhealthy after it
// is cancelled.
//
// Check observes c.source, not ctx. The ctx parameter belongs to the evaluator
// call and controls the evaluation attempt; c.source represents the component's
// shutdown or drain state.
func (c sourceContextCheck) Check(ctx context.Context) Result {
	select {
	case <-c.source.Done():
		res := Unhealthy(c.name, c.reason, c.message)

		if cause := context.Cause(c.source); cause != nil {
			res = res.WithCause(cause)
		}

		return res

	default:
		return Healthy(c.name)
	}
}
