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

func TestParseVersionCanonicalValues(t *testing.T) {
	requireCanonicalSet(t, []string{
		"v1",
		"v2",
		"v10",
		"v1alpha1",
		"v1alpha10",
		"v1beta1",
		"v10beta2",
	}, ParseVersion)
}

func TestParseVersionRejectsNonCanonicalValues(t *testing.T) {
	requireRejectedSet(t, []string{
		"",
		"v0",
		"v0alpha1",
		"v01",
		"v1alpha0",
		"v1alpha01",
		"v1beta0",
		"v1rc1",
		"v1gamma1",
		"v1Alpha1",
		"V1",
		"1",
		"v",
		" v1",
		"v1 ",
		"v1 alpha1",
		"v1alpha1extra",
		"v1alpha1beta1",
		"v1.0.0",
	}, ParseVersion)
}
