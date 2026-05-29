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

// Package types defines closed structural descriptors for ARCORIS API values.
//
// The package owns the descriptor algebra only. A Type is the normalized
// descriptor / IR returned by builders. It can describe primitive values,
// exact numeric widths, exact temporal kinds, object fields, list semantics,
// map values, and references to named structural definitions. FieldDescriptor
// describes object member name, presence, description, and value type. A
// TypeDefinition binds a TypeName to a reusable descriptor. These concepts are
// intentionally structural: they describe shape and descriptor constraints, not
// concrete object storage or runtime behavior.
//
// The field-first builder DSL is the primary construction surface. Field starts
// an object-field declaration, and the following typed method chooses the field
// value descriptor. Unnamed reusable type builders use the same closed TypeExpr
// path without field state.
//
// A typical descriptor declaration keeps object shape, field presence,
// nullability, and type-specific constraints in one readable chain:
//
//	nameType := String()
//	nameType = nameType.MinLen(1)
//	nameType = nameType.MaxLen(253)
//	nameType = nameType.Pattern("^[a-z][a-z0-9-]*$")
//
//	conditionType := Object(
//		Field("type").String().
//			Required().
//			MinLen(1).
//			Description("Stable condition type."),
//
//		Field("status").String().
//			Required().
//			Enum("True", "False", "Unknown"),
//
//		Field("message").String().
//			Optional().
//			Nullable().
//			MaxLen(1024),
//	).UnknownFields(UnknownReject)
//
//	workloadType := Object(
//		Field("spec").Object(
//			Field("name").Ref("arcoris.meta.Name").
//				Required().
//				Description("API name resolved through a type catalog."),
//
//			Field("replicas").Int32().
//				Optional().
//				Range(1, 1000),
//
//			Field("labels").MapOf(
//				String().
//					MinLen(1),
//			).
//				Optional().
//				MaxLen(64),
//		).
//			Required().
//			UnknownFields(UnknownReject),
//
//		Field("status").Object(
//			Field("conditions").ListOf(
//				Ref("arcoris.meta.Condition"),
//			).
//				Optional().
//				Map("type"),
//		).
//			Optional().
//			UnknownFields(UnknownReject),
//	).UnknownFields(UnknownReject)
//
//	nameDef := Define("arcoris.meta.Name", nameType)
//	conditionDef := Define("arcoris.meta.Condition", conditionType)
//	_ = ValidateDefinition(nameDef, nil)
//	_ = ValidateDefinition(conditionDef, nil)
//	_ = ValidateType(workloadType.Type(), nil)
//
// Constructors and builders are intentionally allocation-light and mostly defer
// diagnostics. They produce descriptor values; ValidateType and
// ValidateDefinition perform descriptor-shape validation and return classified
// errors with descriptor paths. The package does not validate concrete object
// values, apply defaults, prune unknown data, export schemas, run codecs, or
// execute arbitrary Go validators.
//
// Type values are immutable-by-convention value objects and normalized
// descriptors. Callers build them through package-owned constructors and
// field-first builders instead of filling public structs. The private exact
// payload-slot layout keeps the type system closed:
// external packages cannot inject arbitrary Go implementations, reflection
// types, runtime objects, callbacks, validators, or transport-specific schema
// fragments. View methods are the read-only inspection API and return detached
// copies of slice-bearing payloads so callers cannot mutate descriptors through
// views, resolved definitions, or field accessors.
//
// Field presence and nullability are separate axes. Required and Optional
// describe whether an object field key must be present. Nullable describes
// whether the value carried by that key may be null. TypeNull is the null
// literal itself; it is not a marker that makes another type nullable.
//
// TypeRef resolves through the Resolver interface defined in this package
// because references are part of the structural type algebra. Concrete mutable
// storage for named definitions is deliberately outside this package; for the
// standard owner-created catalog, use package arcoris.dev/apimachinery/api/typecatalog.
// This separation keeps structural descriptors independent from catalog
// storage, resource registries, runtime schemes, codecs, converters, and global
// registration state.
//
// Recursive TypeDefinition graphs are not supported by api/types. Recursive
// schemas require a future explicit design pass because recursion affects
// validation, export, code generation, and value traversal semantics.
//
// This package deliberately does not implement a value model, casts, coercion,
// defaults, examples as descriptor data, resource definitions, metadata, object
// runtime machinery, JSON or OpenAPI export, concrete value validation,
// pruning, patch/apply, field ownership, Go reflection support, arbitrary
// custom validators, concrete catalogs, or global registries. Those are
// separate design passes above this foundational descriptor layer.
package types
