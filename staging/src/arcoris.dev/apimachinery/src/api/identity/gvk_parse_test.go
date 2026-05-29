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

func TestParseGroupVersionKindCanonicalValues(t *testing.T) {
	requireCanonicalSet(t, []string{
		"v1#Pod",
		"control.arcoris.dev/v1#Worker",
	}, ParseGroupVersionKind)
}

func TestParseGroupVersionKindRejectsLegacyValues(t *testing.T) {
	requireRejectedSet(t, []string{
		"",
		"apps/v1, Kind=Pod",
		"apps/v1, Resource=pods",
		"Kind.version.group",
		"v1#pod",
		"v1",
	}, ParseGroupVersionKind)
}
