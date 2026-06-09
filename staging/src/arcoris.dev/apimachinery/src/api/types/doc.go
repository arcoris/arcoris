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
// The package owns descriptor algebra only. A Descriptor is the normalized
// descriptor / IR returned by builders. It can describe primitive values, exact
// numeric widths, exact temporal kinds, object fields, list semantics, map
// values, and references to named structural definitions. FieldDescriptor
// describes object member name, presence, description, and value descriptor.
// Definition binds a TypeName to a reusable descriptor. These concepts are
// intentionally structural: they describe shape and descriptor constraints, not
// concrete object storage or runtime behavior.
//
// # Model
//
// The implementation follows one exact descriptor path:
//
//	DescriptorKind -> exact payload -> exact view -> exact validator
//
// Public builders are the construction API. They hold a descriptorHeader plus their
// own exact local payload. Calling Descriptor() normalizes that builder state into
// the closed Descriptor IR:
//
//	<Value>Descriptor builder -> Descriptor
//
// Descriptor is then consumed by field descriptors, descriptor definitions, resolvers,
// catalogs, validators, and future exporters. Descriptor is not a builder workspace,
// a value model, a runtime object, or a generic extension point.
//
// Public views are the inspection API:
//
//	Descriptor -> <Value>View
//
// View methods check the DescriptorKind first and return ok=false for the wrong
// descriptor kind. When the DescriptorKind matches, views expose detached read-only
// payload data. Slice-bearing views and nested Descriptor-returning views clone data
// before returning it, so callers cannot mutate descriptors through inspection.
//
// # Code Orientation
//
// The package is organized around descriptor responsibilities rather than one
// large central file. Use these naming patterns when navigating the code:
//
//   - descriptor.go owns the normalized Descriptor layout and the core methods shared by
//     every descriptor kind: Code, IsZero, IsValid, and Nullable.
//   - descriptor_kind.go owns DescriptorKind names and broad classification helpers.
//   - descriptor_header.go and descriptor_flags.go own private builder-header mechanics
//     shared by exact builders.
//   - files named for a descriptor kind own that kind's public builder surface
//     and domain comments. Sibling field files own the matching field wrapper
//     used after Field("name").
//   - *_payload.go files own exact private payload structs plus clone/empty
//     helpers for that payload slot.
//   - *_view.go files own Descriptor.As<Value>() accessors and read-only view methods
//     such as Min, Max, Enum, Fields, Element, Value, and Name.
//   - *_validate.go files own exact validation rules for the matching
//     descriptor kind.
//   - field_* files own field names, presence, finalized FieldDescriptor
//     values, FieldBuilder, sealed FieldExpr mechanics, and field state.
//   - definition.go, type_name.go, resolver.go, and ref files own named
//     descriptor references and the minimal Resolver contract.
//   - descriptor_error.go owns structured validation diagnostics. Broad sentinel
//     errors remain stable for errors.Is, while DescriptorError Reason and Detail
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
//			Field("name").Ref("meta.arcoris.dev.Name").
//				Required().
//				Description("Stable API name resolved through a type catalog."),
//
//			Field("replicas").Int32().
//				Optional().
//				Range(1, 1000),
//
//			Field("image").String().
//				Required().
//				MinBytes(1).
//				MaxBytes(512).
//				Pattern("^[^\\s]+$").
//				Description("Container image reference."),
//
//			Field("labels").MapOf(
//				String().
//					MinBytes(1).
//					MaxBytes(63),
//			).
//				Optional().
//				MaxEntries(64),
//		).
//			Required().
//			UnknownFields(UnknownReject),
//
//		Field("status").Object(
//			Field("conditions").ListOf(
//				Ref("meta.arcoris.dev.Condition"),
//			).
//				Optional().
//				Map("type"),
//		).
//			Optional().
//			UnknownFields(UnknownReject),
//	).UnknownFields(UnknownReject)
//
// # Reusable Unnamed Descriptor Builders
//
// The same builder API can be used without FieldBuilder when a structural descriptor
// should be reused before it is named, referenced, or embedded:
//
//	nameType := String().
//		MinBytes(1).
//		MaxBytes(253).
//		Pattern("^[a-z][a-z0-9-]*$")
//
//	replicasType := Int32().
//		Min(1).
//		Max(1000)
//
//	tagListType := ListOf(
//		String().
//			MinBytes(1).
//			MaxBytes(63),
//	).Set().MaxItems(32)
//
//	_ = nameType
//	_ = replicasType
//	_ = tagListType
//
// # Named Definitions And References
//
// Definition gives a reusable descriptor a TypeName. Ref stores the name in
// a DescriptorRef payload and resolves it through a Resolver during validation.
// ValidateLocal performs descriptor-local structural validation and checks
// reference-name syntax without resolving references. ValidateResolved requires
// a non-nil Resolver and checks that references exist and resolved definitions
// are valid:
//
//	nameDef := Define(
//		"meta.arcoris.dev.Name",
//		String().
//			MinBytes(1).
//			MaxBytes(253),
//	)
//
//	conditionDef := Define(
//		"meta.arcoris.dev.Condition",
//		Object(
//			Field("type").String().
//				Required().
//				MinBytes(1),
//			Field("status").String().
//				Required().
//				Enum("True", "False", "Unknown"),
//		).UnknownFields(UnknownReject),
//	)
//
//	_ = ValidateDefinitionLocal(nameDef)
//	_ = ValidateDefinitionLocal(conditionDef)
//	_ = ValidateLocal(Ref("meta.arcoris.dev.Name").Descriptor())
//
// Recursive Definition graphs are not supported by api/types. Recursive
// schemas require a future explicit design pass because recursion affects
// validation, export, code generation, and value traversal semantics.
//
// # Inspecting Descriptors
//
// Descriptor views are the safe read-only way to inspect normalized descriptors. The
// ok result is part of the API and should be checked before using a view:
//
//	desc := String().
//		Nullable().
//		MinBytes(1).
//		MaxBytes(253).
//		Enum("default", "system").
//		Descriptor()
//
//	if view, ok := desc.AsString(); ok {
//		min, hasMin := view.MinBytes()
//		max, hasMax := view.MaxBytes()
//		values := view.Enum()
//
//		_ = min
//		_ = hasMin
//		_ = max
//		_ = hasMax
//		_ = values
//	}
//
//	if _, ok := desc.AsObject(); !ok {
//		// The descriptor is not an object. No object payload is exposed.
//	}
//
// Inspect object, list, map, and ref descriptors through their exact views:
//
//	objectDescriptor := Object(
//		Field("name").String().
//			Required().
//			MinBytes(1),
//	).Descriptor()
//
//	if view, ok := objectDescriptor.AsObject(); ok {
//		fields := view.Fields()
//		unknown := view.UnknownFields()
//
//		_ = fields
//		_ = unknown
//	}
//
// # Validating Descriptors
//
// Builders mostly defer diagnostics. ValidateLocal checks descriptor-local
// shape; ValidateResolved and ValidateDefinitionResolved additionally require a
// non-nil Resolver and reject unresolved references. Use errors.Is for broad
// programmatic classification and errors.As for precise path, reason, and
// detail:
//
//	err := ValidateLocal(
//		Float64().
//			Min(10).
//			Max(1).
//			Descriptor(),
//	)
//
//	var descriptorErr *DescriptorError
//	if errors.As(err, &descriptorErr) {
//		path := descriptorErr.Path
//		reason := descriptorErr.Reason
//		detail := descriptorErr.Detail
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
//			MaxBytes(1024),
//
//		Field("name").String().
//			Required().
//			MinBytes(1),
//	)
//
// DescriptorNull is the null literal itself. It is not a marker that makes another
// descriptor nullable, and it cannot be marked nullable.
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
