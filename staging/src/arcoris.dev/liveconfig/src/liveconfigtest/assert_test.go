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
	"testing"

	"arcoris.dev/snapshot"
)

func TestAssertionHelpersAcceptValidSnapshots(t *testing.T) {
	snap := snapshot.Snapshot[Config]{
		Revision: snapshot.ZeroRevision.Next(),
		Value:    NewConfig(),
	}

	RequireNonZeroRevision(t, snap)
	RequireRevision(t, snap, snapshot.ZeroRevision.Next())
	RequireChangedSince(t, snap, snapshot.ZeroRevision)
	RequireUnchangedSince(t, snap, snap.Revision)
	RequireSnapshotValue(t, snap, NewConfig(), EqualConfig)

	stamped := snapshot.Stamped[Config]{
		Revision: snap.Revision,
		Value:    snap.Value,
	}
	RequireStampedNonZeroRevision(t, stamped)
	RequireStampedValue(t, stamped, NewConfig(), EqualConfig)
	RequireConfigStampedValue(t, stamped, NewConfig())

	src := NewControlledSource(NewConfig())
	RequireSourceRevision(t, src, snapshot.ZeroRevision.Next())
	RequireSourceValue(t, src, NewConfig(), EqualConfig)
	RequireConfigSourceValue(t, src, NewConfig())
}
