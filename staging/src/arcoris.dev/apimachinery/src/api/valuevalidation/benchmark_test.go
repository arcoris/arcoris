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

package valuevalidation_test

import (
	"fmt"
	"testing"

	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
	"arcoris.dev/apimachinery/api/valuevalidation"
)

func BenchmarkValidateScalarString(b *testing.B) {
	benchmarkValidate(b, value.StringValue("api"), types.String().MinBytes(1).Pattern(`^[a-z]+$`).Descriptor(), valuevalidation.Options{})
}

func BenchmarkValidateScalarInteger(b *testing.B) {
	benchmarkValidate(b, value.Int64Value(42), types.Int64().Min(0).Max(100).Descriptor(), valuevalidation.Options{})
}

func BenchmarkValidateRecordSmall(b *testing.B) {
	payload, descriptor := benchmarkRecordValue(8)

	benchmarkValidate(b, payload, descriptor, valuevalidation.Options{})
}

func BenchmarkValidateRecordLarge(b *testing.B) {
	payload, descriptor := benchmarkRecordValue(100)

	benchmarkValidate(b, payload, descriptor, valuevalidation.Options{})
}

func BenchmarkValidateMapLarge(b *testing.B) {
	payload := benchmarkStringRecord(100)
	descriptor := types.MapOf(types.String().MinBytes(1)).Descriptor()

	benchmarkValidate(b, payload, descriptor, valuevalidation.Options{})
}

func BenchmarkValidateOrderedListLarge(b *testing.B) {
	payload := benchmarkStringList(100)
	descriptor := types.ListOf(types.String().MinBytes(1)).Ordered().Descriptor()

	benchmarkValidate(b, payload, descriptor, valuevalidation.Options{})
}

func BenchmarkValidateListMapLarge(b *testing.B) {
	payload := benchmarkConditionList(100)
	descriptor := types.ListOf(
		types.Object(
			types.Field("type").String().Required(),
			types.Field("status").String().Required(),
		),
	).Map("type").Descriptor()

	benchmarkValidate(b, payload, descriptor, valuevalidation.Options{})
}

func BenchmarkValidateListSetLarge(b *testing.B) {
	payload := benchmarkStringList(100)
	descriptor := types.ListOf(types.String()).Set().Descriptor()

	benchmarkValidate(b, payload, descriptor, valuevalidation.Options{})
}

func BenchmarkValidateStringPatternCached(b *testing.B) {
	payload := benchmarkStringList(100)
	descriptor := types.ListOf(types.String().Pattern(`^[a-z0-9-]+$`)).Ordered().Descriptor()

	benchmarkValidate(b, payload, descriptor, valuevalidation.Options{})
}

func BenchmarkValidateMaxErrors(b *testing.B) {
	fields := make([]types.FieldExpr, 100)
	for i := range fields {
		fields[i] = types.Field(fmt.Sprintf("f%03d", i)).String().Required()
	}
	descriptor := types.Object(fields...).Descriptor()
	payload := value.MustRecordValue()

	benchmarkValidate(b, payload, descriptor, valuevalidation.Options{MaxErrors: 10})
}

func benchmarkValidate(b *testing.B, payload value.Value, descriptor types.Descriptor, opts valuevalidation.Options) {
	b.Helper()
	validator := valuevalidation.New(opts)

	for i := 0; i < b.N; i++ {
		_ = validator.Validate(payload, descriptor)
	}
}

func benchmarkRecordValue(count int) (value.Value, types.Descriptor) {
	members := make([]value.RecordMember, 0, count)
	fields := make([]types.FieldExpr, 0, count)
	for i := 0; i < count; i++ {
		name := fmt.Sprintf("f%03d", i)
		members = append(members, value.MustRecordMember(name, value.StringValue("value")))
		fields = append(fields, types.Field(name).String().Required())
	}

	return value.MustRecordValue(members...), types.Object(fields...).Descriptor()
}

func benchmarkStringRecord(count int) value.Value {
	members := make([]value.RecordMember, 0, count)
	for i := 0; i < count; i++ {
		name := fmt.Sprintf("k%03d", i)
		members = append(members, value.MustRecordMember(name, value.StringValue("value")))
	}

	return value.MustRecordValue(members...)
}

func benchmarkStringList(count int) value.Value {
	items := make([]value.Value, 0, count)
	for i := 0; i < count; i++ {
		items = append(items, value.StringValue(fmt.Sprintf("v%03d", i)))
	}

	return value.MustListValue(items...)
}

func benchmarkConditionList(count int) value.Value {
	items := make([]value.Value, 0, count)
	for i := 0; i < count; i++ {
		items = append(items, value.MustRecordValue(
			value.MustRecordMember("type", value.StringValue(fmt.Sprintf("Ready-%03d", i))),
			value.MustRecordMember("status", value.StringValue("True")),
		))
	}

	return value.MustListValue(items...)
}
