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

import (
	"testing"

	"arcoris.dev/admission"
	panicassert "arcoris.dev/testutil/panic"
)

func TestCatalogNilReceiverPanics(t *testing.T) {
	t.Parallel()

	var catalog *Catalog
	panicassert.RequireMessage(t, nilCatalogPanic, func() {
		_, _ = catalog.Reason(admission.ReasonDenied)
	})
	panicassert.RequireMessage(t, nilCatalogPanic, func() {
		_, _ = catalog.Kind(testKindBulkhead)
	})
	panicassert.RequireMessage(t, nilCatalogPanic, func() {
		_, _ = catalog.Component("resilience.bulkhead")
	})
	panicassert.RequireMessage(t, nilCatalogPanic, func() {
		_ = catalog.Reasons()
	})
	panicassert.RequireMessage(t, nilCatalogPanic, func() {
		_ = catalog.Kinds()
	})
	panicassert.RequireMessage(t, nilCatalogPanic, func() {
		_ = catalog.Components()
	})
	panicassert.RequireMessage(t, nilCatalogPanic, func() {
		_ = catalog.RegisterReason(testReasonDescriptor("custom_reason"))
	})
	panicassert.RequireMessage(t, nilCatalogPanic, func() {
		_ = catalog.RegisterKind(testKindDescriptor("custom_kind"))
	})
	panicassert.RequireMessage(t, nilCatalogPanic, func() {
		_ = catalog.RegisterComponent(testComponentDescriptor("custom.component", testKindBulkhead))
	})
	panicassert.RequireMessage(t, nilCatalogPanic, func() {
		_ = catalog.LenReasons()
	})
	panicassert.RequireMessage(t, nilCatalogPanic, func() {
		_ = catalog.LenKinds()
	})
	panicassert.RequireMessage(t, nilCatalogPanic, func() {
		_ = catalog.LenComponents()
	})
}
