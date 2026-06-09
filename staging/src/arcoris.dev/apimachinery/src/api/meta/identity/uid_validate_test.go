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

package identity

import (
	"strings"
	"testing"
)

func TestUIDValidate(t *testing.T) {
	tests := []string{
		"uid-1",
		"abc.DEF_123",
		"tenant:object",
		strings.Repeat("a", 128),
	}
	for _, value := range tests {
		t.Run("valid/"+value, func(t *testing.T) {
			requireNoError(t, UID(value).Validate())
		})
	}

	for _, value := range []string{
		"",
		"uid/1",
		"uid 1",
		"uid\n1",
		"uid@1",
		strings.Repeat("a", 129),
	} {
		t.Run("invalid/"+value, func(t *testing.T) {
			requireErrorIs(t, UID(value).Validate(), ErrInvalidUID)
		})
	}
}
