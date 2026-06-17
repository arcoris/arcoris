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

import metaidentity "arcoris.dev/apimachinery/api/meta/identity"

// NamespaceRequirementKind classifies whether a predicate has a safe namespace
// requirement that lifecycle/storage adapters may inspect.
type NamespaceRequirementKind uint8

// Namespace requirement classifications.
const (
	// NamespaceUnconstrained means the query imposes no namespace requirement.
	NamespaceUnconstrained NamespaceRequirementKind = iota
	// NamespaceSingle means every match must have Namespace.
	NamespaceSingle
	// NamespaceContradictory means the query contains incompatible namespace
	// requirements and cannot match any object.
	NamespaceContradictory
	// NamespaceDisjunctive means the query permits multiple namespaces.
	NamespaceDisjunctive
	// NamespaceResidual means NOT or another shape prevents safe extraction.
	NamespaceResidual
)

// NamespaceRequirement is a conservative namespace summary of a Predicate.
type NamespaceRequirement struct {
	// Kind identifies the extracted namespace relationship.
	Kind NamespaceRequirementKind
	// Namespace is populated only when Kind is NamespaceSingle.
	Namespace metaidentity.Namespace
}

// NamespaceRequirement returns a conservative namespace summary for p.
func (p Predicate) NamespaceRequirement() NamespaceRequirement {
	return namespaceRequirementForExpr(p.expr)
}

// namespaceRequirementForExpr recursively extracts namespace constraints
// without expanding boolean expressions.
func namespaceRequirementForExpr(e *expr) NamespaceRequirement {
	if e == nil {
		return NamespaceRequirement{Kind: NamespaceUnconstrained}
	}

	switch e.kind {
	case exprNone:
		return NamespaceRequirement{Kind: NamespaceContradictory}
	case exprTerm:
		return namespaceRequirementForTerm(e.term)
	case exprAnd:
		return namespaceRequirementForAnd(e.children)
	case exprOr:
		return namespaceRequirementForOr(e.children)
	case exprNot:
		return NamespaceRequirement{Kind: NamespaceResidual}
	default:
		return NamespaceRequirement{Kind: NamespaceResidual}
	}
}

// namespaceRequirementForTerm extracts namespace facts from key identity terms.
func namespaceRequirementForTerm(t term) NamespaceRequirement {
	switch t.kind {
	case termNamespace, termObject:
		return NamespaceRequirement{Kind: NamespaceSingle, Namespace: t.namespace}
	case termKey:
		return NamespaceRequirement{Kind: NamespaceSingle, Namespace: t.key.Object.Namespace}
	default:
		return NamespaceRequirement{Kind: NamespaceUnconstrained}
	}
}

// namespaceRequirementForAnd combines conjunctive requirements.
func namespaceRequirementForAnd(children []*expr) NamespaceRequirement {
	out := NamespaceRequirement{Kind: NamespaceUnconstrained}
	for _, child := range children {
		next := namespaceRequirementForExpr(child)
		out = mergeNamespaceAnd(out, next)
		if out.Kind == NamespaceContradictory {
			return out
		}
	}

	return out
}

// mergeNamespaceAnd applies AND semantics to namespace summaries.
func mergeNamespaceAnd(left NamespaceRequirement, right NamespaceRequirement) NamespaceRequirement {
	if left.Kind == NamespaceContradictory || right.Kind == NamespaceContradictory {
		return NamespaceRequirement{Kind: NamespaceContradictory}
	}
	if left.Kind == NamespaceUnconstrained {
		return right
	}
	if right.Kind == NamespaceUnconstrained {
		return left
	}
	if left.Kind == NamespaceSingle && right.Kind == NamespaceSingle {
		if left.Namespace == right.Namespace {
			return left
		}
		return NamespaceRequirement{Kind: NamespaceContradictory}
	}
	if left.Kind == NamespaceSingle {
		return left
	}
	if right.Kind == NamespaceSingle {
		return right
	}

	return NamespaceRequirement{Kind: NamespaceResidual}
}

// namespaceRequirementForOr combines disjunctive requirements conservatively.
func namespaceRequirementForOr(children []*expr) NamespaceRequirement {
	out := NamespaceRequirement{Kind: NamespaceContradictory}
	for _, child := range children {
		next := namespaceRequirementForExpr(child)
		out = mergeNamespaceOr(out, next)
		if out.Kind == NamespaceResidual || out.Kind == NamespaceUnconstrained {
			return out
		}
	}

	return out
}

// mergeNamespaceOr applies OR semantics to namespace summaries.
func mergeNamespaceOr(left NamespaceRequirement, right NamespaceRequirement) NamespaceRequirement {
	if left.Kind == NamespaceContradictory {
		return right
	}
	if right.Kind == NamespaceContradictory {
		return left
	}
	if left.Kind == NamespaceUnconstrained || right.Kind == NamespaceUnconstrained {
		return NamespaceRequirement{Kind: NamespaceUnconstrained}
	}
	if left.Kind == NamespaceResidual || right.Kind == NamespaceResidual {
		return NamespaceRequirement{Kind: NamespaceResidual}
	}
	if left.Kind == NamespaceDisjunctive || right.Kind == NamespaceDisjunctive {
		return NamespaceRequirement{Kind: NamespaceDisjunctive}
	}
	if left.Kind == NamespaceSingle && right.Kind == NamespaceSingle {
		if left.Namespace == right.Namespace {
			return left
		}
		return NamespaceRequirement{Kind: NamespaceDisjunctive}
	}

	return NamespaceRequirement{Kind: NamespaceResidual}
}
