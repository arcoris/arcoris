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

// Compile validates and canonicalizes query into an immutable Predicate.
//
// Compile is the only validation path. Callers that only need validation
// should still use Query.Validate, which delegates here.
func Compile(query Query, opts ...CompileOption) (Predicate, error) {
	options := newCompileOptions(opts)
	normalized := normalizeExpr(query.expr)
	compiled, err := compileExpr(normalized, options)
	if err != nil {
		return Predicate{}, invalidQueryError(err)
	}
	compiled.plan = buildPlan(compiled.expr)

	return compiled, nil
}

// compileExpr validates an already-normalized expression and wraps it in a
// Predicate. Plan extraction is intentionally done by Compile after validation.
func compileExpr(e *expr, opts compileOptions) (Predicate, error) {
	canonical, err := validateExpr(e, opts)
	if err != nil {
		return Predicate{}, err
	}

	return Predicate{expr: canonical}, nil
}
