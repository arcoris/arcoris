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

package stamp

import (
	"encoding/json"
	"testing"
	"time"
)

func TestDeletion(t *testing.T) {
	if !(Deletion{}).IsZero() {
		t.Fatal("zero deletion is not zero")
	}

	deletion := Deletion{
		DeletedAt: NewTimestamp(
			time.Date(2026, 5, 30, 12, 0, 0, 0, time.UTC),
		),
		GracePeriodSeconds: 30,
	}
	if deletion.IsZero() {
		t.Fatal("non-zero deletion is zero")
	}
}

func TestDeletionJSONFields(t *testing.T) {
	deletion := Deletion{
		DeletedAt: NewTimestamp(
			time.Date(2026, 5, 30, 12, 0, 0, 0, time.UTC),
		),
		GracePeriodSeconds: 30,
	}

	data, err := json.Marshal(deletion)
	requireNoError(t, err)

	var got map[string]any
	requireNoError(t, json.Unmarshal(data, &got))

	if got["deletedAt"] != "2026-05-30T12:00:00Z" {
		t.Fatalf("deletedAt = %#v", got["deletedAt"])
	}
	if got["gracePeriodSeconds"] != float64(30) {
		t.Fatalf("gracePeriodSeconds = %#v", got["gracePeriodSeconds"])
	}
	if _, ok := got["DeletedAt"]; ok {
		t.Fatalf("unexpected Go field name in JSON: %s", data)
	}
}
