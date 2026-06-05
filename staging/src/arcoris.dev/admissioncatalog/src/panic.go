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

// nilBuilderPanic is the stable panic string for nil *Builder receivers.
//
// A nil builder is a setup bug. Ordinary invalid descriptor input is returned
// as an error by declaration and build methods.
const nilBuilderPanic = "admissioncatalog.Builder: nil builder"

// nilCatalogPanic is the stable panic string for nil *Catalog receivers.
//
// Missing catalog entries are lookup misses. A nil catalog pointer means the
// owner failed to wire a catalog at all, so methods panic close to the bug.
const nilCatalogPanic = "admissioncatalog.Catalog: nil catalog"
