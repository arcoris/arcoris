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

package contextstop

import (
	"context"
	"errors"
)

// Cause returns the most specific cause available for a completed context.
//
// err should usually be ctx.Err. The explicit err parameter keeps callers in
// control of the context-stop boundary and provides a fallback for custom
// Context implementations. When context.Cause(ctx) exposes a custom cause,
// Cause joins it with err so errors.Is can still match context.Canceled or
// context.DeadlineExceeded while also matching the custom cause.
func Cause(ctx context.Context, err error) error {
	cause := context.Cause(ctx)
	if cause == nil {
		return err
	}
	if err == nil || errors.Is(cause, err) {
		return cause
	}

	return errors.Join(err, cause)
}
