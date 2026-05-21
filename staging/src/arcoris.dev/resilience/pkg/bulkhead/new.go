/*
  Copyright 2026 The ARCORIS Authors

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
*/

package bulkhead

import "arcoris.dev/capacity"

// New creates a Bulkhead with the provided local in-flight limit.
//
// A zero limit is valid. It creates a closed bulkhead that rejects acquisition
// until SetLimit raises the limit. This follows capacity.Ledger semantics and
// keeps "closed for now" distinct from an invalid object.
func New(limit Amount) *Bulkhead {
	return &Bulkhead{
		ledger: capacity.NewLedger(capacity.Amount(limit)),
	}
}
