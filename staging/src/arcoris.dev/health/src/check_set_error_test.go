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

import (
	"errors"
	"strings"
	"testing"
)

func TestDuplicateCheckNameErrorClassificationAndFields(t *testing.T) {
	t.Parallel()

	err := DuplicateCheckNameError{
		Name:          "storage",
		Index:         2,
		PreviousIndex: 0,
	}

	if !errors.Is(err, ErrDuplicateCheckName) {
		t.Fatal("DuplicateCheckNameError should match ErrDuplicateCheckName")
	}
	if !strings.Contains(err.Error(), `name="storage"`) {
		t.Fatalf("Error() = %q, want check name", err.Error())
	}
}
