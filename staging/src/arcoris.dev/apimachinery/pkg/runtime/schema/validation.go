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
func invalid(name, val, reason string) error {
	return fmt.Errorf("schema: invalid %s %q: %s", name, val, reason)
}

// invalidValue wraps an error from a nested identifier while preserving the
// outer composite identity that failed. This keeps errors useful when, for
// example, a GroupVersionResource is rejected because only its version segment
// is malformed.
func invalidValue(name, val string, err error) error {
	return fmt.Errorf("schema: invalid %s %q: %w", name, val, err)
}

// nilUnmarshalReceiver reports programmer misuse of an Unmarshal method. It is
// returned instead of panicking so schema remains predictable at API boundaries.
func nilUnmarshalReceiver(name string) error {
	return fmt.Errorf("schema: cannot unmarshal %s into nil receiver", name)
}

// marshalJSONString serializes the already canonical String form after running
// the type's full Validate contract. This rejects invalid direct literals even
// when a caller constructed a value without using the strict parser.
func marshalJSONString(name, val string, validate func() error) ([]byte, error) {
	if err := validate(); err != nil {
		return nil, err
	}
	return json.Marshal(val)
}

// unmarshalJSONString enforces the scalar JSON contract for every identifier.
// Object form, arrays, numbers, booleans, and null are rejected before the
// caller-specific strict parser sees the string payload.
func unmarshalJSONString(name string, data []byte) (string, error) {
	var val *string
	if err := json.Unmarshal(data, &val); err != nil {
		return "", fmt.Errorf("schema: invalid %s JSON: expected string: %w", name, err)
	}
	if val == nil {
		return "", invalid(name+" JSON", "null", "expected string")
	}
	return *val, nil
}

// validateGroupValue enforces the API group grammar. The empty value is valid
// and means the core group; every non-empty value must be a DNS-1123 subdomain
// made of one or more DNS labels separated by single dots.
func validateGroupValue(val string) error {
	if val == "" {
		return nil
	}
	if len(val) > maxDNS1123SubdomainLength {
		return invalid("group", val, "group length must be no more than 253 characters")
	}
	labels := strings.Split(val, ".")
	for _, label := range labels {
		if label == "" {
			return invalid("group", val, "group must not contain empty DNS labels")
		}
		if err := validateDNS1123LabelValue(label); err != nil {
			return invalid("group", val, "group must be a DNS-1123 subdomain")
		}
	}
	return nil
}

// validateDNS1123SingleLabel validates resource-like identifiers that must be a
// single DNS label. Resource requires a non-empty label; Subresource sets
// allowEmpty because an empty subresource is the canonical "absent" value.
func validateDNS1123SingleLabel(name, val string, allowEmpty bool) error {
	if val == "" {
		if allowEmpty {
			return nil
		}
		return invalid(name, val, name+" must be non-empty")
	}
	if strings.Contains(val, ".") {
		return invalid(name, val, name+" must be a DNS-1123 single label")
	}
	if err := validateDNS1123LabelValue(val); err != nil {
		return invalid(name, val, name+" must be a DNS-1123 single label")
	}
	return nil
}

// validateDNS1123LabelValue validates one DNS-1123 label without interpreting
// dots as separators. It is shared by group labels, resources, and
// subresources so their edge-character and length behavior remains identical.
func validateDNS1123LabelValue(val string) error {
	if val == "" {
		return fmt.Errorf("DNS label must be non-empty")
	}
	if len(val) > maxDNS1123LabelLength {
		return fmt.Errorf("DNS label length must be no more than 63 characters")
	}
	if !isDNS1123LabelEdge(val[0]) || !isDNS1123LabelEdge(val[len(val)-1]) {
		return fmt.Errorf("DNS label must start and end with a lowercase letter or digit")
	}
	for i := 1; i < len(val)-1; i++ {
		if !isDNS1123LabelChar(val[i]) {
			return fmt.Errorf("DNS label may contain only lowercase letters, digits, and '-'")
		}
	}
	return nil
}

// validateVersionValue validates the strict ARCORIS version token grammar:
// vN, vNalphaM, or vNbetaM. The implementation is deliberately manual rather
// than regexp-based so leading-zero and suffix rules stay explicit.
func validateVersionValue(val string) error {
	if val == "" {
		return invalid("version", val, "version must be non-empty")
	}
	if len(val) < 2 || val[0] != 'v' {
		return invalid("version", val, "version must match vN, vNalphaM, or vNbetaM")
	}

	i := 1
	if val[i] == '0' {
		i++
	} else if isASCIINonZeroDigit(val[i]) {
		i++
		for i < len(val) && isASCIIDigit(val[i]) {
			i++
		}
	} else {
		return invalid("version", val, "version must match vN, vNalphaM, or vNbetaM")
	}

	if i == len(val) {
		return nil
	}

	if strings.HasPrefix(val[i:], "alpha") {
		i += len("alpha")
	} else if strings.HasPrefix(val[i:], "beta") {
		i += len("beta")
	} else {
		return invalid("version", val, "version must match vN, vNalphaM, or vNbetaM")
	}

	if i == len(val) || !isASCIINonZeroDigit(val[i]) {
		return invalid("version", val, "version must match vN, vNalphaM, or vNbetaM")
	}
	i++
	for i < len(val) && isASCIIDigit(val[i]) {
		i++
	}
	if i != len(val) {
		return invalid("version", val, "version must match vN, vNalphaM, or vNbetaM")
	}
	return nil
}

// validateKindValue validates the schema kind grammar. Kinds are ASCII-only
// because they are API identifiers, not user display strings, and must be
// stable across serializers and generated code.
func validateKindValue(val string) error {
	if val == "" {
		return invalid("kind", val, "kind must be non-empty")
	}
	if !isASCIIUpper(val[0]) {
		return invalid("kind", val, "kind must start with an uppercase ASCII letter")
	}
	for i := 1; i < len(val); i++ {
		if !isASCIIAlpha(val[i]) && !isASCIIDigit(val[i]) {
			return invalid("kind", val, "kind may contain only ASCII letters and digits")
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
