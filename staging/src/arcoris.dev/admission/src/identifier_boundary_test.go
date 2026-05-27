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

import (
	"strings"
	"testing"
)

func TestReasonValidationBoundaries(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		reason Reason
		want   bool
	}{
		{name: "empty", reason: "", want: false},
		{name: "max length", reason: Reason("a" + strings.Repeat("b", maxReasonLength-1)), want: true},
		{name: "max length plus one", reason: Reason("a" + strings.Repeat("b", maxReasonLength)), want: false},
		{name: "starts with digit", reason: "1reason", want: false},
		{name: "starts with underscore", reason: "_reason", want: false},
		{name: "trailing underscore", reason: "reason_", want: false},
		{name: "double underscore", reason: "reason__code", want: false},
		{name: "uppercase", reason: "Reason", want: false},
		{name: "hyphen", reason: "reason-code", want: false},
		{name: "slash", reason: "reason/code", want: false},
		{name: "whitespace", reason: "reason code", want: false},
		{name: "non ascii", reason: "réason", want: false},
		{name: "digits after first", reason: "reason_123", want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.reason.IsValid(); got != tt.want {
				t.Fatalf("%q IsValid() = %t, want %t", tt.reason, got, tt.want)
			}
		})
	}
}

func TestComponentIDValidationBoundaries(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		id   ComponentID
		want bool
	}{
		{name: "empty", id: "", want: false},
		{name: "max length", id: ComponentID("a" + strings.Repeat("b", maxComponentIDLength-1)), want: true},
		{name: "max length plus one", id: ComponentID("a" + strings.Repeat("b", maxComponentIDLength)), want: false},
		{name: "starts with digit", id: "1component", want: false},
		{name: "starts with underscore", id: "_component", want: false},
		{name: "trailing underscore", id: "component_", want: false},
		{name: "double underscore", id: "component__id", want: false},
		{name: "uppercase", id: "Component", want: false},
		{name: "hyphen", id: "component-id", want: false},
		{name: "slash", id: "component/id", want: false},
		{name: "whitespace", id: "component id", want: false},
		{name: "non ascii", id: "compönent", want: false},
		{name: "valid multi segment", id: "resilience.bulkhead_v2", want: true},
		{name: "dot path empty segment", id: "resilience..bulkhead", want: false},
		{name: "dot path leading dot", id: ".resilience", want: false},
		{name: "dot path trailing dot", id: "resilience.", want: false},
		{name: "dot path double dot", id: "resilience..retry", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.id.IsValid(); got != tt.want {
				t.Fatalf("%q IsValid() = %t, want %t", tt.id, got, tt.want)
			}
		})
	}
}

func TestComponentKindValidationBoundaries(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		kind ComponentKind
		want bool
	}{
		{name: "empty", kind: "", want: false},
		{name: "max length", kind: ComponentKind("a" + strings.Repeat("b", maxComponentKindLength-1)), want: true},
		{name: "max length plus one", kind: ComponentKind("a" + strings.Repeat("b", maxComponentKindLength)), want: false},
		{name: "starts with digit", kind: "1kind", want: false},
		{name: "starts with underscore", kind: "_kind", want: false},
		{name: "trailing underscore", kind: "kind_", want: false},
		{name: "double underscore", kind: "kind__name", want: false},
		{name: "uppercase", kind: "Kind", want: false},
		{name: "hyphen", kind: "kind-name", want: false},
		{name: "slash", kind: "kind/name", want: false},
		{name: "whitespace", kind: "kind name", want: false},
		{name: "non ascii", kind: "kïnd", want: false},
		{name: "digits after first", kind: "kind_v2", want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.kind.IsValid(); got != tt.want {
				t.Fatalf("%q IsValid() = %t, want %t", tt.kind, got, tt.want)
			}
		})
	}
}

func TestIdentifierValidationRejectsDynamicOrUnstableSyntax(t *testing.T) {
	t.Parallel()

	invalid := []struct {
		name string
		ok   bool
	}{
		{name: "reason with request path", ok: Reason("request/123").IsValid()},
		{name: "reason with raw error", ok: Reason("deadline exceeded").IsValid()},
		{name: "component with transport address", ok: ComponentID("host.local:8080").IsValid()},
		{name: "component with uuid punctuation", ok: ComponentID("tenant.123e4567-e89b").IsValid()},
		{name: "kind with namespace slash", ok: ComponentKind("tenant/kind").IsValid()},
	}

	for _, tt := range invalid {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if tt.ok {
				t.Fatal("identifier accepted unstable syntax")
			}
		})
	}

	// Syntax validation cannot prove whether a stable-looking value embeds
	// dynamic data. Caller-owned catalogs must govern semantic stability.
	if !Reason("tenant_123").IsValid() {
		t.Fatal("lower_snake_case with digits after the first character should remain syntactically valid")
	}
}
