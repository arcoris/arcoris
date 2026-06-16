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

// TestCanonicalStringsSortsDeduplicatesAndDetaches verifies metadata value-set
// canonicalization does not retain caller slices.
func TestCanonicalStringsSortsDeduplicatesAndDetaches(t *testing.T) {
	source := []string{"qa", "prod", "qa"}
	got := canonicalStrings(source)
	source[1] = "dev"

	if len(got) != 2 || got[0] != "prod" || got[1] != "qa" {
		t.Fatalf("canonicalStrings = %#v; want [prod qa]", got)
	}
}

// TestMetadataDomainName verifies diagnostic names stay stable.
func TestMetadataDomainName(t *testing.T) {
	if metadataDomainName(metadataLabels) != "label" {
		t.Fatal("label domain name changed")
	}
	if metadataDomainName(metadataAnnotations) != "annotation" {
		t.Fatal("annotation domain name changed")
	}
	if metadataDomainName(0) != "metadata" {
		t.Fatal("unknown domain name changed")
	}
}
