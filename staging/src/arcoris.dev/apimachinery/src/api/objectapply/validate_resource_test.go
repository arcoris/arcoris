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
	"arcoris.dev/apimachinery/api/types"
)

func TestValidateResource(t *testing.T) {
	err := New(Options{}).validateResource(testRequest().Resource)

	requireNoError(t, err)
}

func TestValidateResourceRejectsInvalidDefinition(t *testing.T) {
	def := resource.NewDefinition(
		apiidentity.Group("control.arcoris.dev"),
		apiidentity.Kind("Worker"),
		apiidentity.Resource("workers"),
		resource.ScopeNamespaced,
	)

	err := New(Options{}).validateResource(def)

	requireErrorIs(t, err, ErrInvalidResource)
	requireErrorIs(t, err, resource.ErrInvalidDefinition)
	requireObjectApplyError(t, err, pathRequestResource, ErrorReasonInvalidResource)
}

func TestValidateResourceRejectsInvalidInputs(t *testing.T) {
	tests := []struct {
		name   string
		def    resource.Definition
		target error
	}{
		{
			name: "invalid group",
			def: resource.NewDefinition(
				apiidentity.Group("bad group"),
				apiidentity.Kind("Worker"),
				apiidentity.Resource("workers"),
				resource.ScopeNamespaced,
				resource.NewVersion("v1", desiredDescriptor(), resource.Exposed(), resource.Canonical()),
			),
		},
		{
			name: "invalid kind",
			def: resource.NewDefinition(
				apiidentity.Group("control.arcoris.dev"),
				apiidentity.Kind("bad kind"),
				apiidentity.Resource("workers"),
				resource.ScopeNamespaced,
				resource.NewVersion("v1", desiredDescriptor(), resource.Exposed(), resource.Canonical()),
			),
		},
		{
			name: "invalid resource",
			def: resource.NewDefinition(
				apiidentity.Group("control.arcoris.dev"),
				apiidentity.Kind("Worker"),
				apiidentity.Resource("bad resource"),
				resource.ScopeNamespaced,
				resource.NewVersion("v1", desiredDescriptor(), resource.Exposed(), resource.Canonical()),
			),
		},
		{
			name: "invalid scope",
			def: resource.NewDefinition(
				apiidentity.Group("control.arcoris.dev"),
				apiidentity.Kind("Worker"),
				apiidentity.Resource("workers"),
				resource.ScopeInvalid,
				resource.NewVersion("v1", desiredDescriptor(), resource.Exposed(), resource.Canonical()),
			),
		},
		{
			name: "no versions",
			def: resource.NewDefinition(
				apiidentity.Group("control.arcoris.dev"),
				apiidentity.Kind("Worker"),
				apiidentity.Resource("workers"),
				resource.ScopeNamespaced,
			),
			target: resource.ErrInvalidDefinition,
		},
		{
			name: "duplicate versions",
			def: resource.NewDefinition(
				apiidentity.Group("control.arcoris.dev"),
				apiidentity.Kind("Worker"),
				apiidentity.Resource("workers"),
				resource.ScopeNamespaced,
				resource.NewVersion("v1", desiredDescriptor(), resource.Exposed(), resource.Canonical()),
				resource.NewVersion("v1", desiredDescriptor()),
			),
			target: resource.ErrInvalidDefinition,
		},
		{
			name: "invalid version",
			def: resource.NewDefinition(
				apiidentity.Group("control.arcoris.dev"),
				apiidentity.Kind("Worker"),
				apiidentity.Resource("workers"),
				resource.ScopeNamespaced,
				resource.NewVersion("bad version", desiredDescriptor(), resource.Exposed(), resource.Canonical()),
			),
		},
		{
			name: "missing desired",
			def: resource.NewDefinition(
				apiidentity.Group("control.arcoris.dev"),
				apiidentity.Kind("Worker"),
				apiidentity.Resource("workers"),
				resource.ScopeNamespaced,
				resource.NewVersion("v1", types.Descriptor{}, resource.Exposed(), resource.Canonical()),
			),
			target: resource.ErrInvalidVersion,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := New(Options{}).validateResource(tt.def)

			requireErrorIs(t, err, ErrInvalidResource)
			if tt.target != nil {
				requireErrorIs(t, err, tt.target)
			}
			requireObjectApplyError(t, err, pathRequestResource, ErrorReasonInvalidResource)
		})
	}
}
