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

import (
	"strings"
	"testing"

	"arcoris.dev/apimachinery/api/value"
)

// TestTermCanonicalKeyDistinguishesDomains verifies canonical term keys encode
// semantic domains as well as operator and values.
func TestTermCanonicalKeyDistinguishesDomains(t *testing.T) {
	label := mustQ(LabelEquals("env", "prod")).expr.term
	annotation := mustQ(AnnotationEquals("env", "prod")).expr.term
	field := mustQ(FieldEquals(fieldRef("spec.image"), value.StringValue("prod"))).expr.term

	if label.canonicalKey() == annotation.canonicalKey() {
		t.Fatalf("label and annotation canonical keys are equal: %q", label.canonicalKey())
	}
	if !strings.HasPrefix(field.canonicalKey(), "field.") {
		t.Fatalf("field canonical key = %q; want field prefix", field.canonicalKey())
	}
}
