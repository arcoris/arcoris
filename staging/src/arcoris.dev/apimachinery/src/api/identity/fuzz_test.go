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

import "testing"

func FuzzParseGroup(f *testing.F) {
	fuzzComparableIdentity[Group, *Group](
		f,
		[]string{"", "control.arcoris.dev", "apps", "Control.arcoris.dev"},
		ParseGroup,
	)
}

func FuzzParseVersion(f *testing.F) {
	fuzzComparableIdentity[Version, *Version](
		f,
		[]string{"v1", "v1alpha1", "v1beta1", "v0", "v01"},
		ParseVersion,
	)
}

func FuzzParseKind(f *testing.F) {
	fuzzComparableIdentity[Kind, *Kind](
		f,
		[]string{"Pod", "HTTPRoute", "pod", "Pod_Status"},
		ParseKind,
	)
}

func FuzzParseResource(f *testing.F) {
	fuzzComparableIdentity[Resource, *Resource](
		f,
		[]string{"pods", "pod-logs", "Pods", "pods/status"},
		ParseResource,
	)
}

func FuzzParseSubresource(f *testing.F) {
	fuzzComparableIdentity[Subresource, *Subresource](
		f,
		[]string{"", "status", "scale", "Status", "status/log"},
		ParseSubresource,
	)
}

func FuzzParseGroupVersion(f *testing.F) {
	fuzzComparableIdentity[GroupVersion, *GroupVersion](
		f,
		[]string{"v1", "control.arcoris.dev/v1"},
		ParseGroupVersion,
	)
}

func FuzzParseGroupKind(f *testing.F) {
	fuzzComparableIdentity[GroupKind, *GroupKind](
		f,
		[]string{"Pod", "control.arcoris.dev#Worker"},
		ParseGroupKind,
	)
}

func FuzzParseGroupResource(f *testing.F) {
	fuzzComparableIdentity[GroupResource, *GroupResource](
		f,
		[]string{"pods", "control.arcoris.dev:workers"},
		ParseGroupResource,
	)
}

func FuzzParseGroupVersionKind(f *testing.F) {
	fuzzComparableIdentity[GroupVersionKind, *GroupVersionKind](
		f,
		[]string{"v1#Pod", "control.arcoris.dev/v1#Worker"},
		ParseGroupVersionKind,
	)
}

func FuzzParseGroupVersionResource(f *testing.F) {
	fuzzComparableIdentity[GroupVersionResource, *GroupVersionResource](
		f,
		[]string{"v1:pods", "control.arcoris.dev/v1:workers"},
		ParseGroupVersionResource,
	)
}

func FuzzParseResourcePath(f *testing.F) {
	fuzzComparableIdentity[ResourcePath, *ResourcePath](
		f,
		[]string{"pods", "pods/status"},
		ParseResourcePath,
	)
}

func FuzzParseGroupVersionResourcePath(f *testing.F) {
	fuzzComparableIdentity[GroupVersionResourcePath, *GroupVersionResourcePath](
		f,
		[]string{"v1:pods", "v1:pods/status", "control.arcoris.dev/v1:workers/status"},
		ParseGroupVersionResourcePath,
	)
}

func fuzzComparableIdentity[T comparableIdentity, PT identityPointer[T]](
	f *testing.F,
	seeds []string,
	parse func(string) (T, error),
) {
	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, input string) {
		value, err := parse(input)
		if err != nil {
			return
		}

		requireNoError(t, value.Validate())
		requireParseStable(t, value, parse)
		requireUnmarshalStable[T, PT](t, value)
	})
}

func requireParseStable[T comparableIdentity](
	t *testing.T,
	value T,
	parse func(string) (T, error),
) {
	t.Helper()

	again, err := parse(value.String())
	requireNoError(t, err)
	if again != value {
		t.Fatalf("parse(String()) = %q, want %q", again.String(), value.String())
	}
}

func requireUnmarshalStable[T comparableIdentity, PT identityPointer[T]](t *testing.T, value T) {
	t.Helper()

	text := PT(new(T))
	requireNoError(t, text.UnmarshalText(mustMarshalText(t, value)))
	if *text != value {
		t.Fatalf("text roundtrip = %q, want %q", text.String(), value.String())
	}

	jsonValue := PT(new(T))
	requireNoError(t, jsonValue.UnmarshalJSON(mustMarshalJSON(t, value)))
	if *jsonValue != value {
		t.Fatalf("JSON roundtrip = %q, want %q", jsonValue.String(), value.String())
	}
}

func mustMarshalText(t *testing.T, value identityMarshaler) []byte {
	t.Helper()
	data, err := value.MarshalText()
	requireNoError(t, err)
	return data
}

func mustMarshalJSON(t *testing.T, value identityMarshaler) []byte {
	t.Helper()
	data, err := value.MarshalJSON()
	requireNoError(t, err)
	return data
}
