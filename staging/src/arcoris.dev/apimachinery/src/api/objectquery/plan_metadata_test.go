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

package objectquery

import "testing"

// TestConstraintsForMetadataTermEmitsPositiveOperatorsOnly verifies labels and
// annotations expose only safe positive narrowing hints.
func TestConstraintsForMetadataTermEmitsPositiveOperatorsOnly(t *testing.T) {
	label := term{kind: termMetadata, metadataDomain: metadataLabels, metadataKey: "env", operator: OperatorEquals, stringValues: []string{"prod"}}
	annotation := term{kind: termMetadata, metadataDomain: metadataAnnotations, metadataKey: "team", operator: OperatorExists}
	negative := term{kind: termMetadata, metadataDomain: metadataLabels, metadataKey: "env", operator: OperatorNotIn, stringValues: []string{"prod"}}

	labelConstraints := constraintsForMetadataTerm(label)
	if len(labelConstraints) != 1 || labelConstraints[0].Kind != ConstraintLabel {
		t.Fatalf("label constraints = %#v", labelConstraints)
	}
	annotationConstraints := constraintsForMetadataTerm(annotation)
	if len(annotationConstraints) != 1 || annotationConstraints[0].Kind != ConstraintAnnotation {
		t.Fatalf("annotation constraints = %#v", annotationConstraints)
	}
	if got := constraintsForMetadataTerm(negative); got != nil {
		t.Fatalf("negative constraints = %#v; want nil", got)
	}
}
