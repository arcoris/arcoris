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

package identity

import "strings"

// Validate checks the strict ARCORIS API version grammar.
//
// The zero value is invalid as a complete version identity. Experimental APIs
// should use alpha or beta forms such as "v1alpha1", not "v0".
func (v Version) Validate() error {
	return validateVersionValue(string(v))
}

// validateVersionValue checks vN, vNalphaM, and vNbetaM forms.
//
// Version validation is deliberately lexical. It does not understand release
// streams, support windows, conversion policy, or preferred-version ordering.
func validateVersionValue(value string) error {
	if value == "" {
		return invalid(identityNameVersion, value, ErrorReasonEmptyValue, detailVersionNonEmpty)
	}

	if len(value) < 2 || value[0] != 'v' {
		return invalid(
			identityNameVersion,
			value,
			ErrorReasonInvalidForm,
			detailVersionCanonicalForm,
		)
	}

	i := 1
	majorStart := i
	if i >= len(value) || !isDigit(value[i]) {
		return invalid(
			identityNameVersion,
			value,
			ErrorReasonInvalidForm,
			detailMajorVersionPositive,
		)
	}

	if value[i] == '0' {
		return invalid(
			identityNameVersion,
			value,
			ErrorReasonInvalidForm,
			detailMajorVersionNoZero,
		)
	}

	for i < len(value) && isDigit(value[i]) {
		i++
	}

	if i-majorStart > 1 && value[majorStart] == '0' {
		return invalid(
			identityNameVersion,
			value,
			ErrorReasonInvalidForm,
			detailMajorVersionNoLeadingZeroes,
		)
	}

	if i == len(value) {
		return nil
	}

	var suffix string
	switch {
	case strings.HasPrefix(value[i:], alphaVersionQualifier):
		suffix = alphaVersionQualifier
	case strings.HasPrefix(value[i:], betaVersionQualifier):
		suffix = betaVersionQualifier
	default:
		return invalid(identityNameVersion, value, ErrorReasonInvalidForm, detailVersionSuffix)
	}

	i += len(suffix)
	if i == len(value) || !isDigit(value[i]) {
		return invalid(
			identityNameVersion,
			value,
			ErrorReasonInvalidForm,
			detailPrereleaseVersionPositive,
		)
	}

	if value[i] == '0' {
		return invalid(
			identityNameVersion,
			value,
			ErrorReasonInvalidForm,
			detailPrereleaseVersionNoZero,
		)
	}

	for i < len(value) && isDigit(value[i]) {
		i++
	}

	if i != len(value) {
		return invalid(
			identityNameVersion,
			value,
			ErrorReasonInvalidCharacter,
			detailVersionASCII,
		)
	}

	return nil
}
