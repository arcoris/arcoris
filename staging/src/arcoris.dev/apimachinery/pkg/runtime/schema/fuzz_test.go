/*
   Copyright 2026 The ARCORIS Authors

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

*/

package schema

import "testing"

// Fuzz tests exercise parser round-trip invariants. They do not replace the
// focused unit tests, but they keep successful parses tied to Validate and the
// canonical String form.

// FuzzParseGroupVersion verifies successful GroupVersion parses round-trip.
func FuzzParseGroupVersion(f *testing.F) {
	for _, seed := range []string{"v1", "control.arcoris.dev/v1alpha1", "", "apps/v1/deployments"} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, input string) {
		parsed, err := ParseGroupVersion(input)
		if err != nil {
			return
		}
		if err := parsed.Validate(); err != nil {
			t.Fatalf("Validate after successful parse returned error: %v", err)
		}
		roundTrip, err := ParseGroupVersion(parsed.String())
		if err != nil {
			t.Fatalf("ParseGroupVersion(String()) returned error: %v", err)
		}
		if roundTrip != parsed {
			t.Fatalf("roundTrip = %+v, want %+v", roundTrip, parsed)
		}
	})
}

// FuzzParseGroupVersionKind verifies successful GVK parses round-trip.
func FuzzParseGroupVersionKind(f *testing.F) {
	for _, seed := range []string{"v1, Kind=Pod", "control.arcoris.dev/v1alpha1, Kind=WorkloadClass", "WorkloadClass.v1alpha1.control.arcoris.dev"} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, input string) {
		parsed, err := ParseGroupVersionKind(input)
		if err != nil {
			return
		}
		if err := parsed.Validate(); err != nil {
			t.Fatalf("Validate after successful parse returned error: %v", err)
		}
		roundTrip, err := ParseGroupVersionKind(parsed.String())
		if err != nil {
			t.Fatalf("ParseGroupVersionKind(String()) returned error: %v", err)
		}
		if roundTrip != parsed {
			t.Fatalf("roundTrip = %+v, want %+v", roundTrip, parsed)
		}
	})
}

// FuzzParseGroupVersionResource verifies successful GVR parses round-trip.
func FuzzParseGroupVersionResource(f *testing.F) {
	for _, seed := range []string{"v1:pods", "control.arcoris.dev/v1alpha1:workloadclasses", "deployments.v1.apps", "apps/v1/deployments"} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, input string) {
		parsed, err := ParseGroupVersionResource(input)
		if err != nil {
			return
		}
		if err := parsed.Validate(); err != nil {
			t.Fatalf("Validate after successful parse returned error: %v", err)
		}
		roundTrip, err := ParseGroupVersionResource(parsed.String())
		if err != nil {
			t.Fatalf("ParseGroupVersionResource(String()) returned error: %v", err)
		}
		if roundTrip != parsed {
			t.Fatalf("roundTrip = %+v, want %+v", roundTrip, parsed)
		}
	})
}

// FuzzParseResourcePath verifies successful ResourcePath parses round-trip.
func FuzzParseResourcePath(f *testing.F) {
	for _, seed := range []string{"pods", "pods/status", "", "pods/status/extra"} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, input string) {
		parsed, err := ParseResourcePath(input)
		if err != nil {
			return
		}
		if err := parsed.Validate(); err != nil {
			t.Fatalf("Validate after successful parse returned error: %v", err)
		}
		roundTrip, err := ParseResourcePath(parsed.String())
		if err != nil {
			t.Fatalf("ParseResourcePath(String()) returned error: %v", err)
		}
		if roundTrip != parsed {
			t.Fatalf("roundTrip = %+v, want %+v", roundTrip, parsed)
		}
	})
}

// FuzzParseGroupVersionResourcePath verifies successful GVRP parses round-trip.
func FuzzParseGroupVersionResourcePath(f *testing.F) {
	for _, seed := range []string{"v1:pods", "v1:pods/status", "control.arcoris.dev/v1alpha1:workloadclasses/status", "apps/v1/deployments/status"} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, input string) {
		parsed, err := ParseGroupVersionResourcePath(input)
		if err != nil {
			return
		}
		if err := parsed.Validate(); err != nil {
			t.Fatalf("Validate after successful parse returned error: %v", err)
		}
		roundTrip, err := ParseGroupVersionResourcePath(parsed.String())
		if err != nil {
			t.Fatalf("ParseGroupVersionResourcePath(String()) returned error: %v", err)
		}
		if roundTrip != parsed {
			t.Fatalf("roundTrip = %+v, want %+v", roundTrip, parsed)
		}
	})
}
