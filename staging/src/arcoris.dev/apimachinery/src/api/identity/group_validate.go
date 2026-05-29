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

// Validate checks the strict group grammar.
//
// The zero value is valid and means the core group. Every non-empty value must
// be a qualified DNS subdomain made of lowercase ASCII labels.
func (g Group) Validate() error {
	return validateGroupValue(string(g))
}

// validateGroupValue checks the strict API group grammar.
//
// Group validation is intentionally stricter than Kubernetes built-in group
// shortcuts: non-core groups must be qualified DNS subdomains and therefore
// contain at least one dot.
func validateGroupValue(value string) error {
	if value == "" {
		return nil
	}
	if len(value) > maxDNS1123SubdomainLength {
		return invalidf(
			identityNameGroup,
			value,
			ErrorReasonInvalidLength,
			"group length must be <= %d bytes",
			maxDNS1123SubdomainLength,
		)
	}

	if !strings.Contains(value, dnsLabelSeparator) {
		return invalid(
			identityNameGroup,
			value,
			ErrorReasonInvalidForm,
			"non-core group must be a qualified DNS subdomain",
		)
	}

	labels := strings.Split(value, dnsLabelSeparator)
	for _, label := range labels {
		if label == "" {
			return invalid(
				identityNameGroup,
				value,
				ErrorReasonInvalidForm,
				"group must not contain empty DNS labels",
			)
		}

		if err := validateDNS1123Label(identityNameGroup, value, label); err != nil {
			return err
		}
	}

	return nil
}
