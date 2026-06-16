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

package objectquery

import (
	"arcoris.dev/apimachinery/api/meta/annotations"
	"arcoris.dev/apimachinery/api/meta/labels"
)

// validateMetadataTerm enforces metadata operator support, key/value lexical
// rules, operator arity, and value-set canonicalization.
func validateMetadataTerm(
	domain metadataDomain,
	op Operator,
	key string,
	values []string,
) ([]string, error) {
	if !metadataOperators.Supports(op) {
		if !op.IsValid() {
			return nil, invalidOperatorError(op)
		}
		return nil, unsupportedOperatorError(op, metadataDomainName(domain))
	}
	if err := validateMetadataKey(domain, key); err != nil {
		return nil, invalidTermError("invalid %s key: %w", metadataDomainName(domain), err)
	}

	switch op {
	case OperatorExists, OperatorDoesNotExist:
		if len(values) != 0 {
			return nil, invalidTermError("%s takes no values", op.String())
		}
		return nil, nil
	case OperatorEquals, OperatorNotEquals:
		if len(values) != 1 {
			return nil, invalidTermError("%s takes exactly one value", op.String())
		}
	case OperatorIn, OperatorNotIn:
		if len(values) == 0 {
			return nil, invalidTermError("%s requires at least one value", op.String())
		}
	}

	for _, value := range values {
		if err := validateMetadataValue(domain, value); err != nil {
			return nil, invalidTermError("invalid %s value: %w", metadataDomainName(domain), err)
		}
	}

	return canonicalStrings(values), nil
}

// validateMetadataKey delegates label and annotation key syntax to their owner
// packages instead of duplicating lexical rules here.
func validateMetadataKey(domain metadataDomain, key string) error {
	switch domain {
	case metadataLabels:
		return labels.Key(key).ValidateLexical()
	case metadataAnnotations:
		return annotations.Key(key).ValidateLexical()
	default:
		return invalidTermError("unknown metadata domain")
	}
}

// validateMetadataValue delegates metadata value syntax to the selected
// metadata domain.
func validateMetadataValue(domain metadataDomain, val string) error {
	switch domain {
	case metadataLabels:
		return labels.Value(val).ValidateLexical()
	case metadataAnnotations:
		return annotations.Value(val).ValidateLexical()
	default:
		return invalidTermError("unknown metadata domain")
	}
}
