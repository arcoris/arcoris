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

package liveconfig

import "arcoris.dev/snapshot"

// New creates a Holder containing initial as its first last-good value.
//
// New applies the same candidate pipeline used by Apply: clone, normalize,
// validate, and publish. If the initial value is rejected, New returns an error
// and no Holder. A successful New always publishes initial at a non-zero
// revision so readers never observe an uninitialized holder.
func New[T any](initial T, opts ...Option[T]) (*Holder[T], error) {
	cfg := newConfig(opts...)
	h := &Holder[T]{
		cfg: cfg,
		pub: snapshot.NewPublisher[T](
			snapshot.WithClock(cfg.clock),
		),
	}

	cur, _, err := h.prepare(initial)
	if err != nil {
		return nil, err
	}

	h.pub.Publish(cur)
	return h, nil
}
