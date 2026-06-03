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

// LenReasons reports the number of registered reasons.
func (c *Catalog) LenReasons() int {
	c.requireNonNil()
	return c.reasons.Len()
}

// LenKinds reports the number of registered component kinds.
func (c *Catalog) LenKinds() int {
	c.requireNonNil()
	return c.kinds.Len()
}

// LenComponents reports the number of registered components.
func (c *Catalog) LenComponents() int {
	c.requireNonNil()
	return c.components.Len()
}
