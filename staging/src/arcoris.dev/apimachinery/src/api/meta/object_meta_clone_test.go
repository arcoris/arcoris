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

package meta

import (
	"testing"

	"arcoris.dev/apimachinery/api/meta/stamp"
)

func TestObjectMetaClone(t *testing.T) {
	deletion := stamp.Deletion{DeletedAt: stamp.NewTimestamp(testTime())}
	meta := validObjectMeta()
	meta.Deletion = &deletion

	cloned := meta.Clone()
	cloned.Labels["role"] = "manager"
	cloned.Annotations["note"] = "changed"
	cloned.OwnerReferences[0].Controller = false
	cloned.Finalizers[0] = "changed"
	cloned.Deletion.GracePeriodSeconds = 99

	if meta.Labels["role"] != "worker" {
		t.Fatal("labels were not detached")
	}
	if meta.Annotations["note"] != "human readable" {
		t.Fatal("annotations were not detached")
	}
	if !meta.OwnerReferences[0].Controller {
		t.Fatal("owner references were not detached")
	}
	if meta.Finalizers[0] != "cleanup" {
		t.Fatal("finalizers were not detached")
	}
	if meta.Deletion.GracePeriodSeconds != 0 {
		t.Fatal("deletion pointer was not detached")
	}
}
