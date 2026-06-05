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

package builtin

import (
	"testing"

	"arcoris.dev/admission"
	"arcoris.dev/admissioncatalog"
)

func requireReason(t *testing.T, descriptors []admissioncatalog.ReasonDescriptor, reason admission.Reason) admissioncatalog.ReasonDescriptor {
	t.Helper()

	for _, descriptor := range descriptors {
		if descriptor.Reason == reason {
			return descriptor
		}
	}
	t.Fatalf("reason %s not found", reason)
	return admissioncatalog.ReasonDescriptor{}
}

func requireKind(t *testing.T, descriptors []admissioncatalog.ComponentKindDescriptor, kind admission.ComponentKind) admissioncatalog.ComponentKindDescriptor {
	t.Helper()

	for _, descriptor := range descriptors {
		if descriptor.Kind == kind {
			return descriptor
		}
	}
	t.Fatalf("kind %s not found", kind)
	return admissioncatalog.ComponentKindDescriptor{}
}

func requireComponent(t *testing.T, descriptors []admissioncatalog.ComponentDescriptor, id admission.ComponentID) admissioncatalog.ComponentDescriptor {
	t.Helper()

	for _, descriptor := range descriptors {
		if descriptor.ID == id {
			return descriptor
		}
	}
	t.Fatalf("component %s not found", id)
	return admissioncatalog.ComponentDescriptor{}
}
