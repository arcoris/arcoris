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

package health

import "testing"

func TestPrepareCheckSetCopiesCallerSliceAndIndexesNames(t *testing.T) {
	t.Parallel()

	storage := mustCheck(t, "storage", Healthy("storage"))
	queue := mustCheck(t, "queue", Healthy("queue"))

	input := []Checker{storage}
	prepared, names, err := prepareCheckSet(input)
	if err != nil {
		t.Fatalf("prepareCheckSet() = %v, want nil", err)
	}

	input[0] = queue
	if prepared[0].Name() != "storage" {
		t.Fatal("prepareCheckSet exposed caller slice mutation")
	}
	if _, ok := names["storage"]; !ok {
		t.Fatal("names index missing storage")
	}
}
