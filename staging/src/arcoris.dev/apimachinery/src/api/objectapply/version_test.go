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

package objectapply

import (
	"testing"

	apiidentity "arcoris.dev/apimachinery/api/identity"
	"arcoris.dev/apimachinery/api/resource"
)

func TestValidateVersionCompatibility(t *testing.T) {
	req := testRequest()

	err := validateVersionCompatibility(req.Live, req.Applied)

	requireNoError(t, err)
}

func TestValidateVersionCompatibilityRejectsMismatch(t *testing.T) {
	req := testRequest()
	req.Applied.TypeMeta = testTypeMeta("v2")

	err := validateVersionCompatibility(req.Live, req.Applied)

	requireErrorIs(t, err, ErrVersionMismatch)
}

func TestSelectVersion(t *testing.T) {
	version, err := selectVersion(testRequest().Live, testRequest().Resource)
	requireNoError(t, err)

	if version.Version() != "v1" {
		t.Fatalf("Version = %q; want v1", version.Version())
	}
}

func TestSelectVersionRejectsMissingVersion(t *testing.T) {
	res := resource.NewDefinition(
		apiidentity.Group("control.arcoris.dev"),
		apiidentity.Kind("Worker"),
		apiidentity.Resource("workers"),
		resource.ScopeNamespaced,
		resource.NewVersion("v2", desiredDescriptor(), resource.Exposed(), resource.Canonical()),
	)

	_, err := selectVersion(testRequest().Live, res)

	requireErrorIs(t, err, ErrInvalidResource)
}
