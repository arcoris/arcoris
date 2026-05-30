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

package object

import "testing"

func TestObjectIdentityProjections(t *testing.T) {
	obj := New[testDesired, testObserved](
		validTypeMeta(),
		validObjectMeta(),
		testDesired{Replicas: 3},
	)

	if got := obj.GroupVersionKind(); got != validTypeMeta().GroupVersionKind() {
		t.Fatalf("GroupVersionKind() = %#v", got)
	}
	if got := obj.ObjectName(); got != validObjectMeta().ObjectName() {
		t.Fatalf("ObjectName() = %#v", got)
	}
	if got := obj.ObjectIdentity(); got != validObjectMeta().ObjectIdentity() {
		t.Fatalf("ObjectIdentity() = %#v", got)
	}
}
