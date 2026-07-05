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

package objectcontrollerwiring

import (
	"arcoris.dev/apimachinery/runtime/objectcache"
	"arcoris.dev/apimachinery/runtime/objectenqueue"
	"arcoris.dev/apimachinery/runtime/objectindex"
	"arcoris.dev/apimachinery/runtime/objectreflector"
	"arcoris.dev/apimachinery/runtime/objectreflectorsink"
)

// newInputFanout installs one input's read-side sinks in the only supported
// order: cache first, then caller-provided indexes, then enqueue.
//
// The returned index slice is detached from config input so accessors cannot be
// affected by later caller mutation. The index objects themselves are not
// cloned; callers intentionally close over those same instances for Lookup.
func newInputFanout(
	cache *objectcache.Cache,
	indexes []*objectindex.Index,
	enqueueSink *objectenqueue.ReflectorSink,
) (*objectreflectorsink.Fanout, []*objectindex.Index, error) {
	copiedIndexes, err := copyIndexes(indexes)
	if err != nil {
		return nil, nil, err
	}

	sinks := make([]objectreflector.Sink, 0, 2+len(copiedIndexes))
	sinks = append(sinks, cache)
	for _, index := range copiedIndexes {
		sinks = append(sinks, index)
	}
	sinks = append(sinks, enqueueSink)

	fanout, err := objectreflectorsink.NewFanout(sinks...)
	if err != nil {
		return nil, nil, err
	}

	return fanout, copiedIndexes, nil
}

// copyIndexes validates configured index sinks and returns a detached slice.
func copyIndexes(indexes []*objectindex.Index) ([]*objectindex.Index, error) {
	if len(indexes) == 0 {
		return nil, nil
	}

	copied := make([]*objectindex.Index, len(indexes))
	for i, index := range indexes {
		if index == nil {
			return nil, ErrNilIndex
		}
		copied[i] = index
	}

	return copied, nil
}
