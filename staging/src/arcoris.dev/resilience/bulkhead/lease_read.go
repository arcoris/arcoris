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

// Amount returns the amount owned by l.
//
// Current Bulkhead acquisition always reserves one unit, but this method keeps
// the ownership observable through the same shape as capacity.Reservation. It is
// valid before and after release.
func (l *Lease) Amount() Amount {
	l.requireReady()
	return Amount(l.reservation.Amount())
}

// Released reports whether l has already been released.
func (l *Lease) Released() bool {
	l.requireReady()
	return l.reservation.Released()
}
