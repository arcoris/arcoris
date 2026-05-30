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
	"time"

	metaidentity "arcoris.dev/apimachinery/api/meta/identity"
)

func testTime() time.Time {
	return time.Date(2026, 5, 30, 12, 0, 0, 0, time.UTC)
}

func TestObjectMeta(t *testing.T) {
	if !(ObjectMeta{}).IsZero() {
		t.Fatal("zero ObjectMeta IsZero() = false")
	}

	meta := validObjectMeta()
	if meta.IsZero() {
		t.Fatal("non-zero ObjectMeta IsZero() = true")
	}

	if meta.ObjectName() != (metaidentity.ObjectName{Namespace: "system", Name: "worker"}) {
		t.Fatalf("ObjectName() = %#v", meta.ObjectName())
	}
	if meta.ObjectIdentity() != (metaidentity.ObjectIdentity{Namespace: "system", Name: "worker", UID: "uid-1"}) {
		t.Fatalf("ObjectIdentity() = %#v", meta.ObjectIdentity())
	}
}
