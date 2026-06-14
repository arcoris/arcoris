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

package objectsurface

import (
	"testing"

	"arcoris.dev/apimachinery/api/apidocument"
)

func TestKindObjectPath(t *testing.T) {
	kinds := Kinds()
	metadata := kinds.Metadata()
	paths := apidocument.Paths().Object()

	tests := []struct {
		name string
		kind Kind
		want apidocument.Path
	}{
		{name: "desired", kind: kinds.Desired(), want: paths.Desired()},
		{name: "observed", kind: kinds.Observed(), want: paths.Observed()},
		{name: "metadata labels", kind: metadata.Labels(), want: paths.Metadata().Labels()},
		{name: "metadata annotations", kind: metadata.Annotations(), want: paths.Metadata().Annotations()},
		{name: "metadata finalizers", kind: metadata.Finalizers(), want: paths.Metadata().Finalizers()},
		{name: "metadata owner references", kind: metadata.OwnerReferences(), want: paths.Metadata().OwnerReferences()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := tt.kind.ObjectPath()
			if !ok {
				t.Fatalf("ObjectPath() ok = false")
			}
			if got != tt.want {
				t.Fatalf("ObjectPath() = %q; want %q", got, tt.want)
			}
		})
	}
}

func TestKindObjectPathRejectsUnknownSurface(t *testing.T) {
	if got, ok := Kind("metadata.name").ObjectPath(); ok || !got.IsZero() {
		t.Fatalf("ObjectPath() = %q, %v; want zero, false", got, ok)
	}
}
