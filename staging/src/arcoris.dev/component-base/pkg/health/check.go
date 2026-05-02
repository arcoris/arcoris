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

package health

import (
	"context"
	"errors"
)

// maxCheckNameLength defines the maximum length of a health check name.
//
// The limit is deliberately small to keep check names stable, low-cardinality,
// and safe to expose through diagnostics and adapters. Dynamic details such as
// resource IDs, addresses, timestamps, and raw error fragments do not belong in
// check names. They belong in reasons, messages, logs, metrics labels, and other dynamic
// details instead.
const maxCheckNameLength = 128

// ErrEmptyCheckName identifies an empty health check name.
//
// Health check names are ownership identifiers used by registries, reports,
// diagnostics, tests, and transport adapters. A checker with an empty name cannot
// be registered safely because it cannot be addressed or distinguished from
// other checks.
var ErrEmptyCheckName = errors.New("health: empty check name")

// ErrInvalidCheckName identifies a health check name that does not satisfy the
// package's stable identifier syntax.
//
// Check names are intentionally restricted so they remain safe for diagnostics,
// metrics labels, logs, reports, tests, and transport adapters. Dynamic details
// such as resource IDs, addresses, timestamps, and raw error fragments do not
// belong in check names.
var ErrInvalidCheckName = errors.New("health: invalid check name")

// Checker produces one health observation for a named component, subsystem, or
// runtime condition.
//
// Checker is intentionally transport-neutral. It does not receive an HTTP
// request, does not return a transport status code, and does not prescribe gRPC,
// metrics, logging, restart, readiness, admission, or scheduler behavior.
// Adapters and higher-level runtime owners interpret checker results according
// to their own target policies.
//
// Name MUST return a stable, non-empty check identifier. The identifier is used
// by registries, reports, diagnostics, tests, and transport adapters. It MUST NOT
// contain secrets, credentials, connection strings, resource IDs, timestamps, raw
// errors, or other dynamic values.
//
// Check observes the current health state and returns a Result. Implementations
// SHOULD observe ctx when evaluation can block, perform I/O, acquire resources,
// wait on other goroutines, or depend on cancellation or deadlines. Pure
// in-memory checks MAY ignore ctx.
//
// Check implementations SHOULD be safe to call repeatedly. They MUST NOT assume
// an exact number of evaluations because evaluators, probes, diagnostics, and
// tests may run them at different cadences.
//
// Check implementations SHOULD return a Result rather than panic. Panic recovery
// is owned by evaluators so one faulty checker cannot crash health aggregation,
// but a checker should still treat panic as a programming error.
type Checker interface {
	Name() string
	Check(ctx context.Context) Result
}

// ValidCheckName reports whether name satisfies the health check identifier
// syntax.
//
// Valid names use lower_snake_case with ASCII lower-case letters, digits, and
// single underscores between name parts. They MUST start with a lower-case
// letter, MUST NOT end with an underscore, MUST NOT contain repeated
// underscores, and MUST NOT exceed the package-defined maximum check name
// length.
//
// The syntax is deliberately close to Reason syntax. Reasons explain why a
// result has a status; check names identify who produced the result. Both should
// remain stable, low-cardinality, and safe to expose through diagnostics and
// adapters.
func ValidCheckName(name string) bool {
	if len(name) == 0 || len(name) > maxCheckNameLength {
		return false
	}

	previousUnderscore := false

	for i := 0; i < len(name); i++ {
		c := name[i]

		switch {
		case c >= 'a' && c <= 'z':
			previousUnderscore = false

		case c >= '0' && c <= '9':
			if i == 0 {
				return false
			}
			previousUnderscore = false

		case c == '_':
			if i == 0 || previousUnderscore {
				return false
			}
			previousUnderscore = true

		default:
			return false
		}
	}

	return !previousUnderscore
}

// ValidateCheckName validates name as a health check identifier.
//
// ValidateCheckName returns ErrEmptyCheckName for an empty name and
// ErrInvalidCheckName for every other invalid name. Callers should use
// errors.Is for classification rather than matching error strings.
func ValidateCheckName(name string) error {
	if name == "" {
		return ErrEmptyCheckName
	}
	if !ValidCheckName(name) {
		return ErrInvalidCheckName
	}

	return nil
}
