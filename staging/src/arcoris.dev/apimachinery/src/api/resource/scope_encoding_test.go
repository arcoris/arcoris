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

package resource

import "testing"

func TestScopeTextEncoding(t *testing.T) {
	text, err := ScopeGlobal.MarshalText()
	requireNoError(t, err)
	requireEqual(t, string(text), scopeTextGlobal)

	var parsed Scope
	requireNoError(t, parsed.UnmarshalText([]byte(scopeTextNamespaced)))
	requireEqual(t, parsed, ScopeNamespaced)
}

func TestScopeJSONEncoding(t *testing.T) {
	requireScopeJSONRoundTrip(t, ScopeGlobal, scopeTextGlobal)
	requireScopeJSONRoundTrip(t, ScopeNamespaced, scopeTextNamespaced)
}

func TestScopeEncodingRejectsInvalidValues(t *testing.T) {
	_, err := ScopeInvalid.MarshalText()
	requireResourceError(t, err, ErrInvalidScope, pathScope, ErrorReasonInvalidScope)

	_, err = ScopeInvalid.MarshalJSON()
	requireResourceError(t, err, ErrInvalidScope, pathScope, ErrorReasonInvalidScope)

	var parsed Scope
	err = parsed.UnmarshalText([]byte("cluster"))
	requireResourceError(t, err, ErrInvalidScope, pathScope, ErrorReasonInvalidScope)
}

func TestScopeJSONRejectsNonStringValues(t *testing.T) {
	for _, payload := range []string{"null", `{}`, `[]`, `1`, `true`} {
		requireScopeJSONRejects(t, payload)
	}
}

func TestScopeJSONRejectsUnsupportedScopeText(t *testing.T) {
	var scope Scope
	err := scope.UnmarshalJSON([]byte(`"cluster"`))
	requireResourceError(t, err, ErrInvalidScope, pathScope, ErrorReasonInvalidScope)
}

func TestScopeEncodingRejectsNilReceiver(t *testing.T) {
	var nilScope *Scope

	err := nilScope.UnmarshalText([]byte(scopeTextGlobal))
	requireResourceError(t, err, ErrNilReceiver, pathScope, ErrorReasonNilReceiver)

	err = nilScope.UnmarshalJSON([]byte(`"` + scopeTextGlobal + `"`))
	requireResourceError(t, err, ErrNilReceiver, pathScope, ErrorReasonNilReceiver)
}
