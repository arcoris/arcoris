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

package objectcache

import (
	"testing"

	"arcoris.dev/apimachinery/api/objectquery"
)

func TestPackageSnapshotListUsesObjectQueryFullScanSemantics(t *testing.T) {
	source := testItems()
	snapshot := mustSnapshot(t, testListResult(20, source...))
	query := objectquery.Query{
		Identity: objectquery.IdentitySelector{
			Namespace: mustNamespaceEquals(t, "system"),
		},
		Labels: mustLabelSelector(
			t,
			mustLabelIn(t, "env", "prod", "qa"),
			mustLabelNotEquals(t, "tier", "frontend"),
		),
		Annotations: mustAnnotationSelector(
			t,
			mustAnnotationNotIn(t, "team", "platform"),
		),
	}

	assertSnapshotListMatchesObjectQueryFullScan(t, snapshot, source, query)
}
