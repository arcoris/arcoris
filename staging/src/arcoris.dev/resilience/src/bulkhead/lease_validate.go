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

// requireReady panics when l is nil or was not returned by a Bulkhead
// acquisition method.
//
// A valid Lease always has an owning ledger and a positive amount. Zero values
// are invalid because they do not represent live capacity ownership and cannot
// be released safely.
func (l *Lease) requireReady() {
	if l == nil || l.ledger == nil || l.amount.IsZero() {
		panic(ErrInvalidLease)
	}
}
