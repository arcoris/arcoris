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

func TestMetadataKindTreeUsesAPIDocumentFields(t *testing.T) {
	metadata := Kinds().Metadata()
	fields := apidocument.Fields().ObjectMeta()

	tests := []struct {
		name string
		got  Kind
		want Kind
	}{
		{name: "labels", got: metadata.Labels(), want: metadataKind(fields.Labels())},
		{name: "annotations", got: metadata.Annotations(), want: metadataKind(fields.Annotations())},
		{name: "finalizers", got: metadata.Finalizers(), want: metadataKind(fields.Finalizers())},
		{name: "owner references", got: metadata.OwnerReferences(), want: metadataKind(fields.OwnerReferences())},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Fatalf("kind = %q; want %q", tt.got, tt.want)
			}
		})
	}
}
