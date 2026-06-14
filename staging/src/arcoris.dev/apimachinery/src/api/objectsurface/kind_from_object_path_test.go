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

func TestKindFromObjectPath(t *testing.T) {
	kinds := Kinds()
	metadata := kinds.Metadata()
	paths := apidocument.Paths().Object()

	tests := []struct {
		name string
		path apidocument.Path
		want Kind
	}{
		{name: "desired", path: paths.Desired(), want: kinds.Desired()},
		{name: "observed", path: paths.Observed(), want: kinds.Observed()},
		{name: "metadata labels", path: paths.Metadata().Labels(), want: metadata.Labels()},
		{name: "metadata annotations", path: paths.Metadata().Annotations(), want: metadata.Annotations()},
		{name: "metadata finalizers", path: paths.Metadata().Finalizers(), want: metadata.Finalizers()},
		{name: "metadata owner references", path: paths.Metadata().OwnerReferences(), want: metadata.OwnerReferences()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := KindFromObjectPath(tt.path)
			if !ok {
				t.Fatalf("KindFromObjectPath() ok = false")
			}
			if got != tt.want {
				t.Fatalf("KindFromObjectPath() = %q; want %q", got, tt.want)
			}
		})
	}
}

func TestKindFromObjectPathRejectsNonSurfacePaths(t *testing.T) {
	tests := []struct {
		name string
		path apidocument.Path
	}{
		{name: "zero", path: ""},
		{name: "type meta", path: apidocument.Paths().TypeMeta().Kind()},
		{name: "object metadata root", path: apidocument.Paths().Object().MetadataPath()},
		{name: "metadata identity field", path: apidocument.Paths().Object().Metadata().Name()},
		{name: "unknown", path: "object.unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := KindFromObjectPath(tt.path)
			if ok || got != "" {
				t.Fatalf("KindFromObjectPath() = %q, %v; want zero, false", got, ok)
			}
		})
	}
}
