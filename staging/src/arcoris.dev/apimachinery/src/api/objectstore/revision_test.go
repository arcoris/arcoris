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

package objectstore

import "testing"

func TestRevisionZeroIsInvalid(t *testing.T) {
	var revision Revision

	if !revision.IsZero() {
		t.Fatalf("zero revision did not report zero")
	}
	if revision.IsValid() {
		t.Fatalf("zero revision reported valid")
	}
}

func TestRevisionCommittedIsValid(t *testing.T) {
	revision := Revision(1)

	if revision.IsZero() {
		t.Fatalf("committed revision reported zero")
	}
	if !revision.IsValid() {
		t.Fatalf("committed revision reported invalid")
	}
}
