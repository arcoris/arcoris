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

package admissioncatalog

import "testing"

func TestBuilderNilReceiverPanics(t *testing.T) {
	var builder *Builder
	requirePanicMessage(t, nilBuilderPanic, func() {
		_ = builder.DeclareReason(reasonDescriptor(testReason))
	})
	requirePanicMessage(t, nilBuilderPanic, func() {
		_ = builder.DeclareKind(kindDescriptor(testKind))
	})
	requirePanicMessage(t, nilBuilderPanic, func() {
		_ = builder.DeclareComponent(componentDescriptor(testComponent, testKind))
	})
	requirePanicMessage(t, nilBuilderPanic, func() {
		_, _ = builder.Build()
	})
	requirePanicMessage(t, nilBuilderPanic, func() {
		_ = builder.MustBuild()
	})
	requirePanicMessage(t, nilBuilderPanic, func() {
		_ = builder.Include(&Catalog{})
	})
}

func TestCatalogNilReceiverPanics(t *testing.T) {
	var catalog *Catalog
	requirePanicMessage(t, nilCatalogPanic, func() {
		_ = catalog.LenReasons()
	})
	requirePanicMessage(t, nilCatalogPanic, func() {
		_, _ = catalog.Reason(testReason)
	})
	requirePanicMessage(t, nilCatalogPanic, func() {
		_ = catalog.HasReason(testReason)
	})
	requirePanicMessage(t, nilCatalogPanic, func() {
		_ = catalog.Reasons()
	})
	requirePanicMessage(t, nilCatalogPanic, func() {
		_ = catalog.LenKinds()
	})
	requirePanicMessage(t, nilCatalogPanic, func() {
		_, _ = catalog.Kind(testKind)
	})
	requirePanicMessage(t, nilCatalogPanic, func() {
		_ = catalog.HasKind(testKind)
	})
	requirePanicMessage(t, nilCatalogPanic, func() {
		_ = catalog.Kinds()
	})
	requirePanicMessage(t, nilCatalogPanic, func() {
		_ = catalog.LenComponents()
	})
	requirePanicMessage(t, nilCatalogPanic, func() {
		_, _ = catalog.Component(testComponent)
	})
	requirePanicMessage(t, nilCatalogPanic, func() {
		_ = catalog.HasComponent(testComponent)
	})
	requirePanicMessage(t, nilCatalogPanic, func() {
		_ = catalog.Components()
	})
}
