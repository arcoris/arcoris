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

import "arcoris.dev/admissioncatalog"

// NewCatalog returns a fresh immutable catalog containing standard admission
// metadata.
//
// A panic means the package's own descriptor literals are internally invalid.
func NewCatalog() *admissioncatalog.Catalog {
	return admissioncatalog.MustBuild(admissioncatalog.Input{
		Reasons:    ReasonDescriptors(),
		Kinds:      KindDescriptors(),
		Components: ComponentDescriptors(),
	})
}
