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

package run

import (
	"context"
	"errors"
)

// Wait blocks until ctx is done and returns the context cancellation cause.
//
// If ctx has no explicit cause, Wait falls back to ctx.Err. Wait panics when ctx
// is nil. Wait is useful as a Task when a Group should include a context-driven
// sentinel task.
func Wait(ctx context.Context) error {
	requireContext(ctx, errNilWaitContext)

	<-ctx.Done()
	return contextCause(ctx)
}

// IsContextStop reports whether err is a standard context stop error.
func IsContextStop(err error) bool {
	return errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded)
}

// IgnoreContextCanceled returns nil when err matches context.Canceled.
func IgnoreContextCanceled(err error) error {
	if errors.Is(err, context.Canceled) {
		return nil
	}
	return err
}

// IgnoreContextStop returns nil when err represents the observed stop of ctx.
//
// The function ignores errors matching context.Cause(ctx), ctx.Err, or one of
// the standard context stop sentinels after ctx has stopped. It preserves
// unrelated errors. IgnoreContextStop panics when ctx is nil.
func IgnoreContextStop(ctx context.Context, err error) error {
	requireContext(ctx, errNilIgnoreContext)

	if err == nil {
		return nil
	}
	if ctx.Err() == nil {
		return err
	}

	if cause := context.Cause(ctx); cause != nil && errors.Is(err, cause) {
		return nil
	}
	if errors.Is(err, ctx.Err()) {
		return nil
	}
	if IsContextStop(err) {
		return nil
	}

	return err
}

// contextCause returns ctx's cancellation cause with a fallback to ctx.Err.
func contextCause(ctx context.Context) error {
	if cause := context.Cause(ctx); cause != nil {
		return cause
	}
	return ctx.Err()
}
