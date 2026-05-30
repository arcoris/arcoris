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

package metagrammar

import "testing"

func TestValidateOpaqueToken(t *testing.T) {
	allowEmpty := OpaqueTokenOptions{AllowEmpty: true, MaxLength: 8}
	requireNonEmpty := OpaqueTokenOptions{AllowEmpty: false, MaxLength: 8}

	if err := ValidateOpaqueToken("token", "", allowEmpty); err != nil {
		t.Fatalf("ValidateOpaqueToken(empty allowed) error = %v", err)
	}
	if err := ValidateOpaqueToken("token", "opaque-1", requireNonEmpty); err != nil {
		t.Fatalf("ValidateOpaqueToken() error = %v", err)
	}
	if err := ValidateOpaqueToken("token", "", requireNonEmpty); err == nil {
		t.Fatal("ValidateOpaqueToken(empty required) error = nil")
	}
	if err := ValidateOpaqueToken("token", "bad token", allowEmpty); err == nil {
		t.Fatal("ValidateOpaqueToken(whitespace) error = nil")
	}
	if err := ValidateOpaqueToken("token", "bad/token", allowEmpty); err == nil {
		t.Fatal("ValidateOpaqueToken(path separator) error = nil")
	}
	if err := ValidateOpaqueToken("token", "bad\ntoken", allowEmpty); err == nil {
		t.Fatal("ValidateOpaqueToken(control byte) error = nil")
	}
	if err := ValidateOpaqueToken("token", "too-long-token", allowEmpty); err == nil {
		t.Fatal("ValidateOpaqueToken(over max length) error = nil")
	}
}
