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
	"fmt"

	"arcoris.dev/apimachinery/api/objectwatch"
)

// processProgress records a source progress boundary without mutating Sink.
//
// Progress is useful for source liveness and future observability, but it is not
// an object mutation. A progress boundary may equal the latest applied change,
// but it must never move behind already applied or already reported progress.
func (r *Reflector) processProgress(event objectwatch.Event) error {
	if event.Revision.Before(r.lastApplied) {
		return nonMonotonicRevisionError(
			fmt.Errorf("progress revision %s is before last applied revision %s", event.Revision, r.lastApplied),
		)
	}
	if event.Revision.Before(r.lastProgress) {
		return nonMonotonicRevisionError(
			fmt.Errorf("progress revision %s is before last progress revision %s", event.Revision, r.lastProgress),
		)
	}
	r.lastProgress = event.Revision

	return nil
}
