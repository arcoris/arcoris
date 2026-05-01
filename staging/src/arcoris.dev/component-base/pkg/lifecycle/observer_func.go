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

// ObserverFunc adapts a function to Observer.
type ObserverFunc func(Transition)

// ObserveLifecycleTransition calls f with transition.
//
// A nil ObserverFunc ignores every transition. This behavior makes optional
// observer wiring safe when an observer is conditionally configured. Production
// code should still avoid intentionally passing nil observer functions because
// doing so usually hides a missing diagnostics integration.
func (f ObserverFunc) ObserveLifecycleTransition(transition Transition) {
	if f == nil {
		return
	}

	f(transition)
}
