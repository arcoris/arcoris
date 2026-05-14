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

package lifecycle

// reduceTransition builds an uncommitted candidate transition.
//
// It is pure: it does not run guards, assign commit metadata, signal waiters, or
// notify observers.
func reduceTransition(from State, event Event, cause error) (transition Transition, ok bool) {
	to, ok := NextState(from, event)
	if !ok {
		return Transition{
			From:  from,
			To:    from,
			Event: event,
			Cause: cause,
		}, false
	}

	return Transition{
		From:  from,
		To:    to,
		Event: event,
		Cause: cause,
	}, true
}
