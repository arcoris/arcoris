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

func TestValidateDNSSubdomain(t *testing.T) {
	for _, value := range []string{"a", "a.b", "control.arcoris.dev", "x.y-z.1a"} {
		t.Run("valid/"+value, func(t *testing.T) {
			if err := ValidateDNSSubdomain(value); err != nil {
				t.Fatalf("ValidateDNSSubdomain() error = %v", err)
			}
		})
	}

	for _, value := range []string{"", ".a", "a.", "a..b", "Control.arcoris.dev", "a/b", "a_b"} {
		t.Run("invalid/"+value, func(t *testing.T) {
			if err := ValidateDNSSubdomain(value); err == nil {
				t.Fatal("ValidateDNSSubdomain() error = nil")
			}
		})
	}
}

func TestValidateQualifiedDNSSubdomain(t *testing.T) {
	if err := ValidateQualifiedDNSSubdomain("control.arcoris.dev"); err != nil {
		t.Fatalf("ValidateQualifiedDNSSubdomain() error = %v", err)
	}

	if err := ValidateQualifiedDNSSubdomain("control"); err == nil {
		t.Fatal("ValidateQualifiedDNSSubdomain() error = nil")
	}
}
