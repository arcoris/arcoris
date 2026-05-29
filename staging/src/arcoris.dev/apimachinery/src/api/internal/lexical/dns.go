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

package lexical

import (
	"fmt"
	"strings"
)

const (
	// MaxDNS1123LabelLength is the DNS single-label byte limit.
	MaxDNS1123LabelLength = 63

	// MaxDNS1123SubdomainLength is the DNS subdomain byte limit.
	MaxDNS1123SubdomainLength = 253
)

// ValidateDNS1123Label validates a lowercase DNS-1123-like single label.
//
// The helper performs lexical validation only. It does not trim, normalize,
// look up DNS records, or accept Unicode identifier characters.
func ValidateDNS1123Label(value string) *Violation {
	if value == "" {
		return violation(ReasonEmptyValue, "DNS label must be non-empty")
	}
	if len(value) > MaxDNS1123LabelLength {
		return violation(
			ReasonInvalidLength,
			fmt.Sprintf("DNS label length must be <= %d bytes", MaxDNS1123LabelLength),
		)
	}
	if !IsDNS1123LabelEdge(value[0]) || !IsDNS1123LabelEdge(value[len(value)-1]) {
		return violation(
			ReasonInvalidEdge,
			"DNS label must start and end with a lowercase letter or digit",
		)
	}
	for i := 0; i < len(value); i++ {
		if !IsDNS1123LabelChar(value[i]) {
			return violation(
				ReasonInvalidCharacter,
				fmt.Sprintf("DNS label contains invalid byte %q", value[i]),
			)
		}
	}
	return nil
}

// ValidateDNS1123Subdomain validates dot-separated DNS-1123-like labels.
func ValidateDNS1123Subdomain(value string) *Violation {
	if value == "" {
		return violation(ReasonEmptyValue, "DNS subdomain must be non-empty")
	}
	if len(value) > MaxDNS1123SubdomainLength {
		return violation(
			ReasonInvalidLength,
			fmt.Sprintf("DNS subdomain length must be <= %d bytes", MaxDNS1123SubdomainLength),
		)
	}

	labels := strings.Split(value, ".")
	for _, label := range labels {
		if label == "" {
			return violation(ReasonInvalidForm, "DNS subdomain must not contain empty DNS labels")
		}
		if err := ValidateDNS1123Label(label); err != nil {
			return err
		}
	}

	return nil
}

// ValidateQualifiedDNS1123Subdomain validates a non-empty subdomain with at
// least two labels.
func ValidateQualifiedDNS1123Subdomain(value string) *Violation {
	if err := ValidateDNS1123Subdomain(value); err != nil {
		return err
	}
	if !strings.Contains(value, ".") {
		return violation(
			ReasonInvalidForm,
			"qualified DNS subdomain must contain at least one dot",
		)
	}
	return nil
}
