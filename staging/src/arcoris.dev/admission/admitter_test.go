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

package admission

import "testing"

func TestAdmitterFunc(t *testing.T) {
	t.Parallel()
	admitter := AdmitterFunc[Unit, NoGrant, NoMetadata](func(Unit) Result[NoGrant, NoMetadata] {
		return AcceptedNoMetadata(ReasonAdmitted)
	})
	var _ Admitter[Unit, NoGrant, NoMetadata] = admitter
	result := admitter.TryAdmit(Unit{})
	if !result.IsValid() {
		t.Fatalf("unexpected invalid result: %+v", result.Decision())
	}
	if !result.IsAdmitted() {
		t.Fatalf("result should be admitted: %+v", result.Decision())
	}
}
