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

package objectstorewatch

import "errors"

const (
	// defaultMaxHistory keeps enough changes for ordinary short reconnects
	// without making the passive wrapper an unbounded in-memory journal.
	defaultMaxHistory = 4096
	// defaultStreamBuffer gives readers a small burst window while preserving
	// the rule that slow streams cannot block writers indefinitely.
	defaultStreamBuffer = 128
)

// Options controls bounded in-memory retention and stream queues.
type Options struct {
	// MaxHistory is the maximum number of committed changes retained for
	// historical watch replay.
	MaxHistory int
	// StreamBuffer is the per-stream event queue capacity.
	StreamBuffer int
}

// Option mutates Store construction options before validation.
type Option func(*Options) error

// DefaultOptions returns the construction defaults used by New.
func DefaultOptions() Options {
	return Options{
		MaxHistory:   defaultMaxHistory,
		StreamBuffer: defaultStreamBuffer,
	}
}

// WithMaxHistory sets the retained committed-change history limit.
//
// The wrapper keeps at most n committed changes for historical replay. Older
// changes are compacted and starts before the compaction boundary fail with
// objectwatch history-unavailable errors.
func WithMaxHistory(n int) Option {
	return func(options *Options) error {
		if n <= 0 {
			return errors.Join(ErrInvalidOption, errors.New("max history must be positive"))
		}
		options.MaxHistory = n
		return nil
	}
}

// WithStreamBuffer sets each watch stream's bounded queue capacity.
//
// A full queue means the stream is too slow to preserve continuity. Writers do
// not wait for space; the stream is terminated with ErrStreamOverflow wrapped
// as objectwatch continuity loss.
func WithStreamBuffer(n int) Option {
	return func(options *Options) error {
		if n <= 0 {
			return errors.Join(ErrInvalidOption, errors.New("stream buffer must be positive"))
		}
		options.StreamBuffer = n
		return nil
	}
}

// applyOptions validates opts and fills defaults.
//
// Nil Option is rejected instead of ignored so accidental option plumbing bugs
// are visible at construction time.
func applyOptions(opts []Option) (Options, error) {
	options := DefaultOptions()
	for _, opt := range opts {
		if opt == nil {
			return Options{}, errors.Join(ErrInvalidOption, errors.New("nil option"))
		}
		if err := opt(&options); err != nil {
			return Options{}, err
		}
	}

	return options, nil
}
