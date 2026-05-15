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

package healthgrpc

const (
	// healthServerUnavailableMessage hides nil adapter internals from RPC clients.
	healthServerUnavailableMessage = "health server unavailable"

	// nilCheckRequestMessage identifies an invalid Check request boundary.
	nilCheckRequestMessage = "health check request is nil"

	// nilListRequestMessage identifies an invalid List request boundary.
	nilListRequestMessage = "health list request is nil"

	// nilWatchRequestMessage identifies an invalid Watch request boundary.
	nilWatchRequestMessage = "health watch request is nil"

	// unknownServiceMessage is the generic Check error for unconfigured services.
	unknownServiceMessage = "unknown service"

	// healthEvaluationFailedMessage hides evaluator error details from Check clients.
	healthEvaluationFailedMessage = "health evaluation failed"

	// tooManyHealthServicesMessage reports the List guardrail without dumping config.
	tooManyHealthServicesMessage = "too many health services"

	// watchEndedMessage normalizes send errors and stream cancellation.
	watchEndedMessage = "health watch stream ended"
)
