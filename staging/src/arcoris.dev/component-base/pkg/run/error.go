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
	"errors"
	"fmt"
)

var (
	// ErrTaskFailed classifies errors produced by failed run tasks.
	//
	// TaskError values match ErrTaskFailed with errors.Is while preserving the
	// original task error through Unwrap.
	ErrTaskFailed = errors.New("run task failed")
)

// TaskError describes a non-nil error returned by a named task.
type TaskError struct {
	// Name is the stable task name supplied to Group.Go.
	Name string

	// Err is the original error returned by the task.
	Err error
}

// Error returns a human-readable task failure message.
func (e TaskError) Error() string {
	if e.Name == "" {
		return fmt.Sprintf("run task failed: %v", e.Err)
	}

	return fmt.Sprintf("run task %q failed: %v", e.Name, e.Err)
}

// Unwrap returns the original task error.
func (e TaskError) Unwrap() error {
	return e.Err
}

// Is reports whether target is ErrTaskFailed.
func (e TaskError) Is(target error) bool {
	return target == ErrTaskFailed
}

// TaskErrors extracts every TaskError contained in err.
//
// The function walks both single-error unwrap chains and multi-error unwrap
// trees produced by errors.Join. Returned errors preserve traversal order. For
// errors returned by Group.Wait in ErrorModeJoin, that order is the
// deterministic task submission order used by the Group. Join-mode Wait results
// are joined even when only one task failed, so callers can use TaskErrors for a
// uniform single- and multi-failure path.
func TaskErrors(err error) []TaskError {
	if err == nil {
		return nil
	}

	var out []TaskError
	collectTaskErrors(err, &out)
	return out
}

// collectTaskErrors appends TaskError values from err's unwrap tree.
func collectTaskErrors(err error, out *[]TaskError) {
	if err == nil {
		return
	}

	switch typed := err.(type) {
	case TaskError:
		*out = append(*out, typed)
	case *TaskError:
		if typed != nil {
			*out = append(*out, *typed)
		}
	}

	if multi, ok := err.(interface{ Unwrap() []error }); ok {
		for _, child := range multi.Unwrap() {
			collectTaskErrors(child, out)
		}
		return
	}

	if single, ok := err.(interface{ Unwrap() error }); ok {
		collectTaskErrors(single.Unwrap(), out)
	}
}
