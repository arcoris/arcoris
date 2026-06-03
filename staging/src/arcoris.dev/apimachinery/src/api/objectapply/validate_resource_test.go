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

func TestValidateResource(t *testing.T) {
	err := newApplier(Options{}).validateResource(testRequest().Resource)

	requireNoError(t, err)
}

func TestValidateResourceRejectsInvalidDefinition(t *testing.T) {
	def := resource.NewDefinition(
		apiidentity.Group("control.arcoris.dev"),
		apiidentity.Kind("Worker"),
		apiidentity.Resource("workers"),
		resource.ScopeNamespaced,
	)

	err := newApplier(Options{}).validateResource(def)

	requireErrorIs(t, err, ErrInvalidResource)
	requireErrorIs(t, err, resource.ErrInvalidDefinition)
	requireObjectApplyError(t, err, pathRequestResource, ErrorReasonInvalidResource)
}
