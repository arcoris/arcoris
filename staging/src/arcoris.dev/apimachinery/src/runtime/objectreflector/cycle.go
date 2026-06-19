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

package objectreflector

import (
	"context"
	"errors"
)

// runCycle performs one ListCollection -> Replace -> Watch -> ApplyChange pass.
//
// A cycle is intentionally linear and unbuffered: each event is read, validated,
// and applied before the next event is requested from the source stream. This
// keeps backpressure visible to the source stream and avoids introducing a
// second queue whose durability and replay semantics would need a separate
// contract.
func (r *Reflector) runCycle(ctx context.Context) (err error) {
	read, err := r.source.ListCollection(ctx, r.collection)
	if err != nil {
		return err
	}
	if err := r.validateCollectionRead(read); err != nil {
		return err
	}
	if err := r.sink.Replace(ctx, read.Clone()); err != nil {
		return err
	}
	r.lastApplied = read.Revision()
	r.lastProgress = read.Revision()

	stream, err := r.openWatch(ctx, read)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := stream.Close(); closeErr != nil {
			err = errors.Join(err, closeErr)
		}
	}()

	return r.consumeStream(ctx, stream)
}
