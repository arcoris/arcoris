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

const (
	// maxDNS1123LabelLength is the DNS single-label byte limit.
	maxDNS1123LabelLength = 63

	// maxDNS1123SubdomainLength is the DNS subdomain byte limit.
	maxDNS1123SubdomainLength = 253
)

// validateDNS1123SingleLabel checks resource-like one-label identities.
func validateDNS1123SingleLabel(name, value string, allowEmpty bool) error {
	if value == "" {
		if allowEmpty {
			return nil
		}
		return invalid(name, value, ErrorReasonEmptyValue, name+" must be non-empty")
	}
	if strings.Contains(value, dnsLabelSeparator) || strings.Contains(value, subresourceSeparator) {
		return invalid(name, value, ErrorReasonInvalidCharacter, name+" must be a DNS-1123 single label")
	}
	return validateDNS1123Label(name, value, value)
}

// validateDNS1123Label checks one DNS-1123 label.
func validateDNS1123Label(name, fullValue, label string) error {
	if len(label) > maxDNS1123LabelLength {
		return invalidf(name, fullValue, ErrorReasonInvalidLength, "DNS label length must be <= %d bytes", maxDNS1123LabelLength)
	}
	if !isDNS1123Edge(label[0]) || !isDNS1123Edge(label[len(label)-1]) {
		return invalid(name, fullValue, ErrorReasonInvalidEdge, "DNS label must start and end with a lowercase letter or digit")
	}
	for i := 0; i < len(label); i++ {
		if !isDNS1123Char(label[i]) {
			return invalidf(name, fullValue, ErrorReasonInvalidCharacter, "DNS label contains invalid byte %q", label[i])
		}
	}
	return nil
}

// isDNS1123Edge reports whether b can start or end a DNS label.
func isDNS1123Edge(b byte) bool {
	return isLower(b) || isDigit(b)
}

// isDNS1123Char reports whether b can appear inside a DNS label.
func isDNS1123Char(b byte) bool {
	return isLower(b) || isDigit(b) || b == '-'
}

// isLower reports whether b is an ASCII lowercase letter.
func isLower(b byte) bool { return b >= 'a' && b <= 'z' }

// isUpper reports whether b is an ASCII uppercase letter.
func isUpper(b byte) bool { return b >= 'A' && b <= 'Z' }

// isDigit reports whether b is an ASCII digit.
func isDigit(b byte) bool { return b >= '0' && b <= '9' }
