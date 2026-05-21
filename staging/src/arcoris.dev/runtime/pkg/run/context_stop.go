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
// The check is specific to the supplied context: err must match
// context.Cause(ctx) or ctx.Err after ctx has stopped. Context errors from
// nested operations are preserved when they do not match this context's own
// stop, even if IsContextStop would classify them broadly. IgnoreContextStop
// panics when ctx is nil.
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

	return err
}
