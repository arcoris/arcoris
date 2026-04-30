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

package schema

import (
	"encoding/json"
	"fmt"
	"strings"
)

// validation.go contains the lexical rules shared by all schema identifiers.
// The helpers in this file are intentionally small and hand-written so the
// package can stay dependency-free and so every rejected byte has one clear
// reason rooted in the ARCORIS schema contract.

const (
	// maxDNS1123LabelLength is the DNS label limit used by groups, resources,
	// and subresources. The package treats names as ASCII wire tokens, so byte
	// length and character length are the same for every valid value.
	maxDNS1123LabelLength = 63
	// maxDNS1123SubdomainLength is the total DNS subdomain limit for named API
	// groups. The empty group bypasses this limit because it is the explicit
	// canonical spelling of the core API group.
	maxDNS1123SubdomainLength = 253
)

// invalid formats a schema error with the public field/type name, the rejected
// value, and a stable human-readable reason. Parse and validation functions
// use this consistently so callers do not need package-specific error types to
// understand failures.
func invalid(name, value, reason string) error {
	return fmt.Errorf("schema: invalid %s %q: %s", name, value, reason)
}

// invalidValue wraps an error from a nested identifier while preserving the
// outer composite identity that failed. This keeps errors useful when, for
// example, a GroupVersionResource is rejected because only its version segment
// is malformed.
func invalidValue(name, value string, err error) error {
	return fmt.Errorf("schema: invalid %s %q: %w", name, value, err)
}

// nilUnmarshalReceiver reports programmer misuse of an Unmarshal method. It is
// returned instead of panicking so schema remains predictable at API boundaries.
func nilUnmarshalReceiver(name string) error {
	return fmt.Errorf("schema: cannot unmarshal %s into nil receiver", name)
}

// marshalJSONString serializes the already canonical String form after running
// the type's full Validate contract. This rejects invalid direct literals even
// when a caller constructed a value without using the strict parser.
func marshalJSONString(name, value string, validate func() error) ([]byte, error) {
	if err := validate(); err != nil {
		return nil, err
	}
	return json.Marshal(value)
}

// unmarshalJSONString enforces the scalar JSON contract for every identifier.
// Object form, arrays, numbers, booleans, and null are rejected before the
// caller-specific strict parser sees the string payload.
func unmarshalJSONString(name string, data []byte) (string, error) {
	var value *string
	if err := json.Unmarshal(data, &value); err != nil {
		return "", fmt.Errorf("schema: invalid %s JSON: expected string: %w", name, err)
	}
	if value == nil {
		return "", invalid(name+" JSON", "null", "expected string")
	}
	return *value, nil
}

// validateGroupValue enforces the API group grammar. The empty value is valid
// and means the core group; every non-empty value must be a DNS-1123 subdomain
// made of one or more DNS labels separated by single dots.
func validateGroupValue(value string) error {
	if value == "" {
		return nil
	}
	if len(value) > maxDNS1123SubdomainLength {
		return invalid("group", value, "group length must be no more than 253 characters")
	}
	labels := strings.Split(value, ".")
	for _, label := range labels {
		if label == "" {
			return invalid("group", value, "group must not contain empty DNS labels")
		}
		if err := validateDNS1123LabelValue(label); err != nil {
			return invalid("group", value, "group must be a DNS-1123 subdomain")
		}
	}
	return nil
}

// validateDNS1123SingleLabel validates resource-like identifiers that must be a
// single DNS label. Resource requires a non-empty label; Subresource sets
// allowEmpty because an empty subresource is the canonical "absent" value.
func validateDNS1123SingleLabel(name, value string, allowEmpty bool) error {
	if value == "" {
		if allowEmpty {
			return nil
		}
		return invalid(name, value, name+" must be non-empty")
	}
	if strings.Contains(value, ".") {
		return invalid(name, value, name+" must be a DNS-1123 single label")
	}
	if err := validateDNS1123LabelValue(value); err != nil {
		return invalid(name, value, name+" must be a DNS-1123 single label")
	}
	return nil
}

// validateDNS1123LabelValue validates one DNS-1123 label without interpreting
// dots as separators. It is shared by group labels, resources, and
// subresources so their edge-character and length behavior remains identical.
func validateDNS1123LabelValue(value string) error {
	if value == "" {
		return fmt.Errorf("DNS label must be non-empty")
	}
	if len(value) > maxDNS1123LabelLength {
		return fmt.Errorf("DNS label length must be no more than 63 characters")
	}
	if !isDNS1123LabelEdge(value[0]) || !isDNS1123LabelEdge(value[len(value)-1]) {
		return fmt.Errorf("DNS label must start and end with a lowercase letter or digit")
	}
	for i := 1; i < len(value)-1; i++ {
		if !isDNS1123LabelChar(value[i]) {
			return fmt.Errorf("DNS label may contain only lowercase letters, digits, and '-'")
		}
	}
	return nil
}

// validateVersionValue validates the strict ARCORIS version token grammar:
// vN, vNalphaM, or vNbetaM. The implementation is deliberately manual rather
// than regexp-based so leading-zero and suffix rules stay explicit.
func validateVersionValue(value string) error {
	if value == "" {
		return invalid("version", value, "version must be non-empty")
	}
	if len(value) < 2 || value[0] != 'v' {
		return invalid("version", value, "version must match vN, vNalphaM, or vNbetaM")
	}

	i := 1
	if value[i] == '0' {
		i++
	} else if isASCIINonZeroDigit(value[i]) {
		i++
		for i < len(value) && isASCIIDigit(value[i]) {
			i++
		}
	} else {
		return invalid("version", value, "version must match vN, vNalphaM, or vNbetaM")
	}

	if i == len(value) {
		return nil
	}

	if strings.HasPrefix(value[i:], "alpha") {
		i += len("alpha")
	} else if strings.HasPrefix(value[i:], "beta") {
		i += len("beta")
	} else {
		return invalid("version", value, "version must match vN, vNalphaM, or vNbetaM")
	}

	if i == len(value) || !isASCIINonZeroDigit(value[i]) {
		return invalid("version", value, "version must match vN, vNalphaM, or vNbetaM")
	}
	i++
	for i < len(value) && isASCIIDigit(value[i]) {
		i++
	}
	if i != len(value) {
		return invalid("version", value, "version must match vN, vNalphaM, or vNbetaM")
	}
	return nil
}

// validateKindValue validates the schema kind grammar. Kinds are ASCII-only
// because they are API identifiers, not user display strings, and must be
// stable across serializers and generated code.
func validateKindValue(value string) error {
	if value == "" {
		return invalid("kind", value, "kind must be non-empty")
	}
	if !isASCIIUpper(value[0]) {
		return invalid("kind", value, "kind must start with an uppercase ASCII letter")
	}
	for i := 1; i < len(value); i++ {
		if !isASCIIAlpha(value[i]) && !isASCIIDigit(value[i]) {
			return invalid("kind", value, "kind may contain only ASCII letters and digits")
		}
	}
	return nil
}

// isDNS1123LabelEdge reports whether a byte can appear at the start or end of
// a DNS-1123 label under the schema contract.
func isDNS1123LabelEdge(ch byte) bool {
	return isASCIILower(ch) || isASCIIDigit(ch)
}

// isDNS1123LabelChar reports whether a byte can appear inside a DNS-1123 label.
func isDNS1123LabelChar(ch byte) bool {
	return isDNS1123LabelEdge(ch) || ch == '-'
}

// isASCIIAlpha reports whether a byte is an ASCII letter.
func isASCIIAlpha(ch byte) bool {
	return isASCIILower(ch) || isASCIIUpper(ch)
}

// isASCIILower reports whether a byte is a lowercase ASCII letter.
func isASCIILower(ch byte) bool {
	return ch >= 'a' && ch <= 'z'
}

// isASCIIUpper reports whether a byte is an uppercase ASCII letter.
func isASCIIUpper(ch byte) bool {
	return ch >= 'A' && ch <= 'Z'
}

// isASCIIDigit reports whether a byte is an ASCII decimal digit.
func isASCIIDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

// isASCIINonZeroDigit reports whether a byte is an ASCII decimal digit that can
// start a non-zero numeric version component.
func isASCIINonZeroDigit(ch byte) bool {
	return ch >= '1' && ch <= '9'
}
