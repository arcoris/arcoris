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
	"sort"
)

// taskErrorRecord stores a TaskError with the task submission sequence.
//
// Submission sequence is separate from completion order. Goroutine completion is
// race-dependent; submission order is deterministic and better for stable
// diagnostics.
type taskErrorRecord struct {
	seq uint64
	err TaskError
}

// recordTaskError records err as a named task failure.
//
// Every task error is appended for Wait, including errors that arrive after
// fail-fast cancellation has already closed the Group for new submissions. When
// cancel-on-error is enabled, only the first recorded task error reserves the
// right to cancel the context. Cancellation is invoked after releasing the Group
// mutex so context cancellation cannot run owner callbacks while submission and
// error state are locked.
func (g *Group) recordTaskError(seq uint64, name string, err error) {
	if err == nil {
		return
	}

	taskErr := TaskError{Name: name, Err: err}
	shouldCancel := false

	g.mu.Lock()
	g.errs = append(g.errs, taskErrorRecord{seq: seq, err: taskErr})
	if g.config.cancelOnError {
		g.closed = true
		if !g.taskCancelSet {
			g.taskCancelSet = true
			shouldCancel = true
		}
	}
	g.mu.Unlock()

	if shouldCancel {
		g.cancel(taskErr)
	}
}

// buildWaitError builds the configured Wait error from recorded task errors.
//
// Wait reports task failures only. Owner and parent cancellation remain context
// causes unless a task chooses to return them as its own error.
// ErrorModeFirst returns the first task error recorded under the Group mutex.
// ErrorModeJoin sorts by submission sequence and intentionally calls
// errors.Join for both single and multiple errors so join-mode callers can rely
// on one result shape.
func (g *Group) buildWaitError() error {
	g.mu.Lock()
	records := append([]taskErrorRecord(nil), g.errs...)
	mode := g.config.errorMode
	g.mu.Unlock()

	if len(records) == 0 {
		return nil
	}

	if mode == ErrorModeFirst {
		return records[0].err
	}

	sort.Slice(records, func(i, j int) bool {
		return records[i].seq < records[j].seq
	})

	errs := make([]error, 0, len(records))
	for _, record := range records {
		errs = append(errs, record.err)
	}

	return errors.Join(errs...)
}
