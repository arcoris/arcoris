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
// The package owns descriptor algebra only. A Type is the normalized
// descriptor / IR returned by builders. It can describe primitive values, exact
// numeric widths, exact temporal kinds, object fields, list semantics, map
// values, and references to named structural definitions. FieldDescriptor
// describes object member name, presence, description, and value type.
// TypeDefinition binds a TypeName to a reusable descriptor. These concepts are
// intentionally structural: they describe shape and descriptor constraints, not
// concrete object storage or runtime behavior.
//
// # Model
//
// The implementation follows one exact descriptor path:
//
//	TypeCode -> exact payload -> exact view -> exact validator
//
// Public builders are the construction API. They hold a typeHeader plus their
// own exact local payload. Calling Type() normalizes that builder state into
// the closed Type IR:
//
//	<Value>Type builder -> Type
//
// Type is then consumed by field descriptors, type definitions, resolvers,
// catalogs, validators, and future exporters. Type is not a builder workspace,
// a value model, a runtime object, or a generic extension point.
//
// Public views are the inspection API:
//
//	Type -> <Value>View
//
// View methods check the TypeCode first and return ok=false for the wrong
// descriptor kind. When the TypeCode matches, views expose detached read-only
// payload data. Slice-bearing views and nested Type-returning views clone data
// before returning it, so callers cannot mutate descriptors through inspection.
//
// # Code Orientation
//
// The package is organized around descriptor responsibilities rather than one
// large central file. Use these naming patterns when navigating the code:
//
//   - type.go owns the normalized Type layout and the core methods shared by
//     every descriptor kind: Code, IsZero, IsValid, and Nullable.
//   - type_code.go owns TypeCode names and broad classification helpers.
//   - type_header.go and type_flags.go own private builder-header mechanics
//     shared by exact builders.
//   - files named for a descriptor kind own that kind's public builder surface
//     and domain comments. Sibling field files own the matching field wrapper
//     used after Field("name").
//   - *_payload.go files own exact private payload structs plus clone/empty
//     helpers for that payload slot.
//   - *_view.go files own Type.<Value>() accessors and read-only view methods
//     such as Min, Max, Enum, Fields, Element, Value, and Name.
//   - *_validate.go files own exact validation rules for the matching
//     descriptor kind.
//   - field_* files own field names, presence, finalized FieldDescriptor
//     values, FieldBuilder, sealed FieldExpr mechanics, and field state.
//   - type_definition.go, type_name.go, resolver.go, and ref files own named
//     type references and the minimal Resolver contract.
//   - type_error.go owns structured validation diagnostics. Broad sentinel
//     errors remain stable for errors.Is, while TypeError Reason and Detail
//     describe the precise descriptor invariant that failed.
//
// Concrete mutable storage for named definitions is intentionally outside this
// package. The owner-created catalog implementation lives in
// arcoris.dev/apimachinery/api/typecatalog and depends on this package, not the
// other way around.
//
// # Field-First Object Declarations
//
// Field starts an object-field declaration. The next method chooses the value
// descriptor kind, and the rest of the chain applies presence, nullability,
// constraints, semantic policy, and optional documentation text. Put
// Description at the end of the chain so structural rules remain visually
// grouped before metadata:
//
//	workloadType := Object(
//		Field("spec").Object(
//			Field("name").Ref("arcoris.meta.Name").
//				Required().
//				Description("Stable API name resolved through a type catalog."),
//
//			Field("replicas").Int32().
//				Optional().
//				Range(1, 1000),
//
//			Field("image").String().
//				Required().
//				MinLen(1).
//				MaxLen(512).
//				Pattern("^[^\\s]+$").
//				Description("Container image reference."),
//
//			Field("labels").MapOf(
//				String().
//					MinLen(1).
//					MaxLen(63),
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
// # Reusable Unnamed Type Builders
//
// The same builder API can be used without FieldBuilder when a structural type
// should be reused before it is named, referenced, or embedded:
//
//	nameType := String().
//		MinLen(1).
//		MaxLen(253).
//		Pattern("^[a-z][a-z0-9-]*$")
//
//	replicasType := Int32().
//		Min(1).
//		Max(1000)
//
//	tagListType := ListOf(
//		String().
//			MinLen(1).
//			MaxLen(63),
//	).Set().MaxLen(32)
//
//	_ = nameType
//	_ = replicasType
//	_ = tagListType
//
// # Named Definitions And References
//
// TypeDefinition gives a reusable descriptor a TypeName. Ref stores the name in
// a TypeRef payload and resolves it through a Resolver during validation. A nil
// resolver performs local structural validation only; a non-nil resolver also
// checks that references exist and that resolved definitions are valid:
//
//	nameDef := Define(
//		"arcoris.meta.Name",
//		String().
//			MinLen(1).
//			MaxLen(253),
//	)
//
//	conditionDef := Define(
//		"arcoris.meta.Condition",
//		Object(
//			Field("type").String().
//				Required().
//				MinLen(1),
//			Field("status").String().
//				Required().
//				Enum("True", "False", "Unknown"),
//		).UnknownFields(UnknownReject),
//	)
//
//	_ = ValidateDefinition(nameDef, nil)
//	_ = ValidateDefinition(conditionDef, nil)
//	_ = ValidateType(Ref("arcoris.meta.Name").Type(), nil)
//
// Recursive TypeDefinition graphs are not supported by api/types. Recursive
// schemas require a future explicit design pass because recursion affects
// validation, export, code generation, and value traversal semantics.
//
// # Inspecting Descriptors
//
// Type views are the safe read-only way to inspect normalized descriptors. The
// ok result is part of the API and should be checked before using a view:
//
//	tp := String().
//		Nullable().
//		MinLen(1).
//		MaxLen(253).
//		Enum("default", "system").
//		Type()
//
//	if view, ok := tp.String(); ok {
//		min, hasMin := view.MinLen()
//		max, hasMax := view.MaxLen()
//		values := view.Enum()
//
//		_ = min
//		_ = hasMin
//		_ = max
//		_ = hasMax
//		_ = values
//	}
//
//	if _, ok := tp.Object(); !ok {
//		// The descriptor is not an object. No object payload is exposed.
//	}
//
// Inspect object, list, map, and ref descriptors through their exact views:
//
//	objectType := Object(
//		Field("name").String().
//			Required().
//			MinLen(1),
//	).Type()
//
//	if view, ok := objectType.Object(); ok {
//		fields := view.Fields()
//		unknown := view.UnknownFields()
//
//		_ = fields
//		_ = unknown
//	}
//
// # Validating Descriptors
//
// Builders mostly defer diagnostics. ValidateType and ValidateDefinition check
// descriptor shape and return structured TypeError diagnostics. Use errors.Is
// for broad programmatic classification and errors.As for precise path, reason,
// and detail:
//
//	err := ValidateType(
//		Float64().
//			Min(10).
//			Max(1).
//			Type(),
//		nil,
//	)
//
//	var typeErr *TypeError
//	if errors.As(err, &typeErr) {
//		path := typeErr.Path
//		reason := typeErr.Reason
//		detail := typeErr.Detail
//
//		_ = path
//		_ = reason
//		_ = detail
//	}
//
// # Presence And Nullability
//
// Required and Optional describe whether an object field key must be present.
// Nullable describes whether the value carried by that key may be null. These
// are separate axes:
//
//	Object(
//		Field("description").String().
//			Optional().
//			Nullable().
//			MaxLen(1024),
//
//		Field("name").String().
//			Required().
//			MinLen(1),
//	)
//
// TypeNull is the null literal itself. It is not a marker that makes another
// type nullable, and it cannot be marked nullable.
//
// # Boundaries
//
// This package deliberately does not implement a value model, casts, coercion,
// defaults, examples as descriptor data, resource definitions, metadata, object
// runtime machinery, JSON or OpenAPI export, concrete value validation,
// pruning, patch/apply, field ownership, Go reflection support, arbitrary
// custom validators, concrete catalogs, or global registries. Those are
// separate design passes above this foundational descriptor layer.
package types
