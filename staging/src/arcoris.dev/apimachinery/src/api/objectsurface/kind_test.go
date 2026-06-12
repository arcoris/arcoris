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

import "testing"

func TestKindClassification(t *testing.T) {
	tests := []struct {
		name     string
		kind     Kind
		text     string
		valid    bool
		ownable  bool
		metadata bool
	}{
		{name: "desired", kind: Desired, text: "desired", valid: true, ownable: true},
		{name: "observed", kind: Observed, text: "observed", valid: true, ownable: true},
		{name: "labels", kind: MetadataLabels, text: "metadata.labels", valid: true, ownable: true, metadata: true},
		{name: "annotations", kind: MetadataAnnotations, text: "metadata.annotations", valid: true, ownable: true, metadata: true},
		{name: "finalizers", kind: MetadataFinalizers, text: "metadata.finalizers", valid: true, metadata: true},
		{name: "owner references", kind: MetadataOwnerReferences, text: "metadata.ownerReferences", valid: true, metadata: true},
		{name: "empty", kind: "", valid: false},
		{name: "unknown", kind: "metadata.name", text: "metadata.name", valid: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.kind.String(); got != tt.text {
				t.Fatalf("String() = %q; want %q", got, tt.text)
			}
			if got := tt.kind.IsValid(); got != tt.valid {
				t.Fatalf("IsValid() = %v; want %v", got, tt.valid)
			}
			if got := tt.kind.IsOwnable(); got != tt.ownable {
				t.Fatalf("IsOwnable() = %v; want %v", got, tt.ownable)
			}
			if got := tt.kind.IsMetadata(); got != tt.metadata {
				t.Fatalf("IsMetadata() = %v; want %v", got, tt.metadata)
			}
		})
	}
}
