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

import "strings"

// QualifiedNameOptions describes an internal "[prefix/]name" grammar.
//
// The helper exists for future metadata-like descriptor keys. It is not a
// public label, annotation, finalizer, or lifecycle name type.
type QualifiedNameOptions struct {
	// RequirePrefix rejects unprefixed names.
	RequirePrefix bool

	// AllowPrefix allows an optional qualified DNS subdomain prefix.
	AllowPrefix bool

	// MaxNameLength bounds the name segment when positive.
	MaxNameLength int

	// AllowNameDot allows "." inside the name segment.
	AllowNameDot bool
}

// ValidateQualifiedName validates an internal "[prefix/]name" token.
func ValidateQualifiedName(value string, opts QualifiedNameOptions) *Violation {
	if value == "" {
		return violation(ReasonEmptyValue, "qualified name must be non-empty")
	}

	if strings.Count(value, "/") > 1 {
		return violation(ReasonInvalidForm, "qualified name must contain at most one slash")
	}

	prefix, name, hasPrefix := strings.Cut(value, "/")
	if !hasPrefix {
		if opts.RequirePrefix {
			return violation(ReasonInvalidForm, "qualified name requires a DNS subdomain prefix")
		}
		return validateQualifiedNameSegment(prefix, opts)
	}

	if !opts.AllowPrefix && !opts.RequirePrefix {
		return violation(ReasonInvalidForm, "qualified name prefix is not allowed")
	}
	if prefix == "" {
		return violation(ReasonInvalidForm, "qualified name prefix must be non-empty")
	}
	if name == "" {
		return violation(ReasonEmptyValue, "qualified name segment must be non-empty")
	}
	if err := ValidateQualifiedDNS1123Subdomain(prefix); err != nil {
		return err
	}

	return validateQualifiedNameSegment(name, opts)
}

// validateQualifiedNameSegment checks the name portion after prefix handling.
func validateQualifiedNameSegment(value string, opts QualifiedNameOptions) *Violation {
	return ValidateASCIIToken(value, TokenOptions{
		MinLength:         1,
		MaxLength:         opts.MaxNameLength,
		AllowLower:        true,
		AllowDigit:        true,
		AllowHyphen:       true,
		AllowDot:          opts.AllowNameDot,
		RequireAlnumEdges: true,
	})
}
