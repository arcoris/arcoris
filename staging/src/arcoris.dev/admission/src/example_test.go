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

package admission_test

import (
	"fmt"

	"arcoris.dev/admission"
	admissionbuiltin "arcoris.dev/admissioncatalog/builtin"
)

type lease struct {
	// id is only here to make the example output visibly read the typed grant.
	id string
}

func ExampleResult() {
	result := admission.Granted(
		admission.ReasonAdmitted,
		&lease{id: "l1"},
		"snapshot-1",
	)
	grant, hasGrant := result.Grant()
	metadata, hasMetadata := result.Metadata()

	fmt.Println(result.IsAdmitted(), result.HasSideEffect())
	fmt.Println(hasGrant, grant.id)
	fmt.Println(hasMetadata, metadata)

	// Output:
	// true true
	// true l1
	// true snapshot-1
}

func ExampleAdmitterFunc() {
	admitter := admission.AdmitterFunc[admission.Unit, admission.NoGrant, admission.NoMetadata](
		func(admission.Unit) admission.Result[admission.NoGrant, admission.NoMetadata] {
			return admission.DeniedNoMetadata(admission.Reason("capacity_exhausted"))
		},
	)

	result := admitter.TryAdmit(admission.Unit{})
	fmt.Println(result.IsDenied(), result.Decision().Reason)

	// Output:
	// true capacity_exhausted
}

func Example_builtinCatalog() {
	catalog := admissionbuiltin.NewCatalog()
	component, ok := catalog.Component("resilience.bulkhead")

	fmt.Println(ok, component.Kind)

	// Output:
	// true bulkhead
}
