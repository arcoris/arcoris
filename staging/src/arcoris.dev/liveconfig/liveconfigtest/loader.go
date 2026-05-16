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

package liveconfigtest

import (
	"context"
	"errors"
	"sync"
)

// ErrNoLoadResult reports that a scripted loader was asked to load more values
// than it was configured to return.
var ErrNoLoadResult = errors.New("liveconfigtest: no scripted load result")

// LoadResult is one scripted result returned by Loader.
//
// A result with Err == nil is a successful load. A result with Err != nil models
// a source-layer failure such as "candidate unavailable" without requiring a
// real file, environment variable, network client, or watcher in the test.
type LoadResult[T any] struct {
	// Value is returned when Err is nil.
	Value T

	// Err is returned as the load error. When Err is non-nil, Value is still
	// stored for test diagnostics but callers should ignore it just as they would
	// ignore a failed production load value.
	Err error
}

// Succeeded reports whether r represents a successful scripted load.
func (r LoadResult[T]) Succeeded() bool {
	return r.Err == nil
}

// Failed reports whether r represents a failed scripted load.
func (r LoadResult[T]) Failed() bool {
	return r.Err != nil
}

// Loaded returns a successful scripted load result.
func Loaded[T any](val T) LoadResult[T] {
	return LoadResult[T]{Value: val}
}

// LoadFailed returns a failed scripted load result.
func LoadFailed[T any](err error) LoadResult[T] {
	return LoadResult[T]{Err: err}
}

// Loader is a deterministic scripted loader for tests.
//
// Each Load call consumes the next configured LoadResult. Loader is safe for
// concurrent use, but results are consumed in call order rather than goroutine
// creation order.
//
// Loader is deliberately not a watcher or retry loop. Tests own the scheduling:
// they call Load when they want a candidate and inspect Calls, Remaining, or
// Exhausted to assert consumption behavior.
type Loader[T any] struct {
	mu      sync.Mutex
	results []LoadResult[T]
	calls   int
}

// NewLoader creates a scripted loader from results.
func NewLoader[T any](results ...LoadResult[T]) *Loader[T] {
	copied := append([]LoadResult[T](nil), results...)
	return &Loader[T]{results: copied}
}

// Append adds results to the end of the script.
//
// Append does not reset the current call index. It is useful for tests that
// model a source producing more candidates after the first reload cycle.
func (l *Loader[T]) Append(results ...LoadResult[T]) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.results = append(l.results, results...)
}

// Reset replaces the script and clears the consumed-call count.
func (l *Loader[T]) Reset(results ...LoadResult[T]) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.results = append(l.results[:0], results...)
	l.calls = 0
}

// Load returns the next scripted result.
//
// If ctx is already done, Load returns ctx.Err without consuming a scripted
// result. If no scripted result remains, Load returns ErrNoLoadResult.
func (l *Loader[T]) Load(ctx context.Context) (T, error) {
	var zero T
	if err := ctx.Err(); err != nil {
		return zero, err
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	if l.calls >= len(l.results) {
		return zero, ErrNoLoadResult
	}

	res := l.results[l.calls]
	l.calls++
	return res.Value, res.Err
}

// Peek returns the next scripted result without consuming it.
func (l *Loader[T]) Peek() (LoadResult[T], bool) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.calls >= len(l.results) {
		var zero LoadResult[T]
		return zero, false
	}
	return l.results[l.calls], true
}

// Calls returns the number of scripted results consumed by Load.
func (l *Loader[T]) Calls() int {
	l.mu.Lock()
	defer l.mu.Unlock()

	return l.calls
}

// Remaining returns the number of scripted results not yet consumed by Load.
func (l *Loader[T]) Remaining() int {
	l.mu.Lock()
	defer l.mu.Unlock()

	return len(l.results) - l.calls
}

// Exhausted reports whether every scripted result has been consumed.
func (l *Loader[T]) Exhausted() bool {
	return l.Remaining() == 0
}
