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

// ObserverChain notifies several observers in order.
//
// ObserverChain is useful when a component wants to attach several independent
// integrations to lifecycle transitions, such as metrics, tracing, debug
// recording, and tests.
//
// Nil observers in the chain are ignored. This makes conditional observer
// composition safe and keeps option-building code simple.
type ObserverChain []Observer

// ObserveLifecycleTransition notifies all observers in order.
//
// The same committed transition value is passed to every observer. The chain
// does not recover panics, does not launch goroutines, and does not provide
// backpressure or retry semantics. Observers that need those behaviors should
// implement them explicitly outside the lifecycle package.
func (o ObserverChain) ObserveLifecycleTransition(transition Transition) {
	for _, observer := range o {
		if observer == nil {
			continue
		}

		observer.ObserveLifecycleTransition(transition)
	}
}

// notifyObservers notifies observers about transition.
//
// This package-local helper is the controller-facing form of ObserverChain. It
// keeps controller code simple while allowing options to store observers as a
// plain slice.
//
// Controller MUST call notifyObservers only after the transition has been
// committed and outside the controller transition lock.
func notifyObservers(observers []Observer, transition Transition) {
	ObserverChain(observers).ObserveLifecycleTransition(transition)
}
