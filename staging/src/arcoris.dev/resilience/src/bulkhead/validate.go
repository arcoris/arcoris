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

package bulkhead

// requireReady panics when b is nil or was not created by New.
//
// Public methods call requireReady before delegating to capacity validation.
// That makes receiver misuse consistently report bulkhead-owned errors even
// when the same call also passes an invalid amount.
func (b *Bulkhead) requireReady() {
	if b == nil {
		panic(ErrNilBulkhead)
	}
	if b.ledger == nil {
		panic(ErrUninitializedBulkhead)
	}
}
