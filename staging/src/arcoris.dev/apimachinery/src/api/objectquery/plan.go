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
	apiidentity "arcoris.dev/apimachinery/api/identity"
	metaidentity "arcoris.dev/apimachinery/api/meta/identity"
	"arcoris.dev/apimachinery/api/objectstore"
	"arcoris.dev/apimachinery/api/value"
)

// ConstraintKind identifies one advisory planning constraint domain.
type ConstraintKind uint8

// Constraint domains emitted by Predicate.Plan.
const (
	// ConstraintResource narrows by objectstore.Key.Resource.
	ConstraintResource ConstraintKind = iota + 1
	// ConstraintObjectNamespace narrows by objectstore.Key.Object.Namespace.
	ConstraintObjectNamespace
	// ConstraintObjectName narrows by objectstore.Key.Object.Name.
	ConstraintObjectName
	// ConstraintObject narrows by namespace/name together.
	ConstraintObject
	// ConstraintKey narrows by the complete objectstore.Key.
	ConstraintKey
	// ConstraintLabel narrows by a positive label requirement.
	ConstraintLabel
	// ConstraintAnnotation narrows by a positive annotation requirement.
	ConstraintAnnotation
	// ConstraintField narrows by a positive registered selectable field term.
	ConstraintField
)

// ConstraintRef identifies what a planning constraint applies to.
type ConstraintRef struct {
	// Resource is populated for ConstraintResource.
	Resource apiidentity.GroupVersionResource
	// Namespace is populated for ConstraintObjectNamespace and ConstraintObject.
	Namespace metaidentity.Namespace
	// Name is populated for ConstraintObjectName and ConstraintObject.
	Name metaidentity.Name
	// Key is populated for ConstraintKey.
	Key objectstore.Key
	// Metadata is the label or annotation key for metadata constraints.
	Metadata string
	// Field is populated for ConstraintField.
	Field FieldRef
}

// Constraint is a positive, safe narrowing hint extracted from a Predicate.
type Constraint struct {
	// Kind identifies the indexed domain that may narrow candidates.
	Kind ConstraintKind
	// Ref identifies the resource, object identity, metadata key, or field.
	Ref ConstraintRef
	// Op is always a positive index-safe operator in the current plan.
	Op Operator
	// Values contains detached literal values when Op carries operands.
	Values []value.Value
}

// Plan contains immutable, conservative narrowing constraints.
type Plan struct {
	// constraints is sorted by the canonical expression tree and contains only
	// safe positive hints. Match remains the semantic source of truth.
	constraints []Constraint
}

// RangeConstraints visits detached constraints in canonical order.
func (p Plan) RangeConstraints(fn func(Constraint) bool) {
	for _, constraint := range p.constraints {
		if !fn(cloneConstraint(constraint)) {
			return
		}
	}
}

// clone returns a detached copy of p so callers cannot mutate Predicate plan
// state through returned constraints.
func (p Plan) clone() Plan {
	if len(p.constraints) == 0 {
		return Plan{}
	}

	out := Plan{constraints: make([]Constraint, len(p.constraints))}
	for i, constraint := range p.constraints {
		out.constraints[i] = cloneConstraint(constraint)
	}

	return out
}

// cloneConstraint detaches literal slices while keeping comparable reference
// values by value.
func cloneConstraint(c Constraint) Constraint {
	if len(c.Values) > 0 {
		c.Values = cloneValues(c.Values)
	}

	return c
}
