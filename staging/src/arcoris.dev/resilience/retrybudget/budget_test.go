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

package retrybudget

import "arcoris.dev/snapshot"

type stubBudget struct {
	snap Snapshot
}

func (b *stubBudget) RecordOriginal() {}

func (b *stubBudget) TryAdmitRetry() Decision {
	return Decision{
		Allowed: true,
		Reason:  ReasonAllowed,
		Snapshot: snapshot.Snapshot[Snapshot]{
			Revision: snapshot.ZeroRevision.Next(),
			Value:    b.snap,
		},
	}
}

func (b *stubBudget) Snapshot() snapshot.Snapshot[Snapshot] {
	return snapshot.Snapshot[Snapshot]{
		Revision: snapshot.ZeroRevision.Next(),
		Value:    b.snap,
	}
}

var _ OriginalRecorder = (*stubBudget)(nil)
var _ RetryAdmitter = (*stubBudget)(nil)
var _ snapshot.Source[Snapshot] = (*stubBudget)(nil)
var _ Budget = (*stubBudget)(nil)
