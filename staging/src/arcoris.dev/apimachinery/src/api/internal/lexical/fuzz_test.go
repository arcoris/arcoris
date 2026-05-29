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

package lexical

import (
	"strings"
	"testing"
)

func FuzzValidateDNS1123Label(f *testing.F) {
	for _, seed := range []string{"workers", "worker-1", "", "Workers", "worker/main", "воркер"} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, value string) {
		if ValidateDNS1123Label(value) != nil {
			return
		}
		if len(value) == 0 || len(value) > MaxDNS1123LabelLength {
			t.Fatalf("valid label has invalid length: %q", value)
		}
		if strings.ContainsAny(value, "./_ ") {
			t.Fatalf("valid label contains forbidden separator: %q", value)
		}
		if !IsDNS1123LabelEdge(value[0]) || !IsDNS1123LabelEdge(value[len(value)-1]) {
			t.Fatalf("valid label has invalid edge: %q", value)
		}
		for i := 0; i < len(value); i++ {
			if !IsDNS1123LabelChar(value[i]) {
				t.Fatalf("valid label contains invalid byte %q: %q", value[i], value)
			}
		}
	})
}

func FuzzValidateDNS1123Subdomain(f *testing.F) {
	for _, seed := range []string{"control.arcoris.dev", "a.b", "", "control..dev", "Control.dev"} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, value string) {
		if ValidateDNS1123Subdomain(value) != nil {
			return
		}
		if len(value) == 0 || len(value) > MaxDNS1123SubdomainLength {
			t.Fatalf("valid subdomain has invalid length: %q", value)
		}
		for _, label := range strings.Split(value, ".") {
			if ValidateDNS1123Label(label) != nil {
				t.Fatalf("valid subdomain has invalid label %q in %q", label, value)
			}
		}
	})
}

func FuzzValidateQualifiedDNS1123Subdomain(f *testing.F) {
	for _, seed := range []string{"control.arcoris.dev", "a.b", "workers", "", "control..dev"} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, value string) {
		if ValidateQualifiedDNS1123Subdomain(value) != nil {
			return
		}
		if !strings.Contains(value, ".") {
			t.Fatalf("valid qualified subdomain lacks dot: %q", value)
		}
		if ValidateDNS1123Subdomain(value) != nil {
			t.Fatalf("valid qualified subdomain is not valid subdomain: %q", value)
		}
	})
}
