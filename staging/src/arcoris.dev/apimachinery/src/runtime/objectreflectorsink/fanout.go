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

package objectreflectorsink

import (
	"context"

	"arcoris.dev/apimachinery/api/objectstore"
	storewatchapi "arcoris.dev/apimachinery/api/objectstorewatch"
	"arcoris.dev/apimachinery/runtime/objectreflector"
)

// Fanout forwards each reflected sink operation to multiple sinks in order.
//
// Fanout is a sequential objectreflector.Sink combinator. It does not clone
// operation payloads, retry failed calls, add atomicity across sinks, or repair
// partial success. The first downstream error stops the operation and is
// returned unchanged.
type Fanout struct {
	// sinks is the immutable ordered sink list installed at construction time.
	sinks []objectreflector.Sink
}

// NewFanout constructs a Sink that forwards Replace and ApplyChange calls to
// sinks in the configured order.
func NewFanout(sinks ...objectreflector.Sink) (*Fanout, error) {
	if len(sinks) == 0 {
		return nil, ErrNoSinks
	}

	ordered := make([]objectreflector.Sink, len(sinks))
	for i, sink := range sinks {
		if isNilSink(sink) {
			return nil, objectreflector.ErrNilSink
		}
		ordered[i] = sink
	}

	return &Fanout{sinks: ordered}, nil
}

// Replace forwards a complete collection read to every configured sink.
//
// Calls are made sequentially in constructor order. Replace stops at the first
// sink error and returns that error unchanged.
func (f *Fanout) Replace(ctx context.Context, read storewatchapi.CollectionRead) error {
	if err := f.validate(); err != nil {
		return err
	}

	for _, sink := range f.sinks {
		if err := sink.Replace(ctx, read); err != nil {
			return err
		}
	}

	return nil
}

// ApplyChange forwards one committed object transition to every configured
// sink.
//
// Calls are made sequentially in constructor order. ApplyChange stops at the
// first sink error and returns that error unchanged.
func (f *Fanout) ApplyChange(ctx context.Context, change objectstore.Change) error {
	if err := f.validate(); err != nil {
		return err
	}

	for _, sink := range f.sinks {
		if err := sink.ApplyChange(ctx, change); err != nil {
			return err
		}
	}

	return nil
}

// validate checks that the receiver still contains a usable ordered sink list.
//
// NewFanout performs the same validation for normal construction. This
// defensive check keeps nil receivers and package-internal corruptions from
// producing delayed panics in reflector control flow.
func (f *Fanout) validate() error {
	if f == nil || len(f.sinks) == 0 {
		return ErrInvalidFanout
	}

	for _, sink := range f.sinks {
		if isNilSink(sink) {
			return ErrInvalidFanout
		}
	}

	return nil
}
