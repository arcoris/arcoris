#!/bin/sh
#
# Copyright 2026 The ARCORIS Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -eu

root=${1:-.}

fail() {
	printf '%s\n' "health boundary violation: $1" >&2
	exit 1
}

has_go_import() {
	dir=$1
	pattern=$2
	grep -R "$pattern" "$root/$dir" --include='*.go' >/dev/null 2>&1
}

has_mod_ref() {
	file=$1
	pattern=$2
	grep "$pattern" "$root/$file" >/dev/null 2>&1
}

health_dir=staging/src/arcoris.dev/health/src
healtheval_dir=staging/src/arcoris.dev/healtheval/src
probe_dir=staging/src/arcoris.dev/healthprobe/src
http_dir=staging/src/arcoris.dev/healthhttp/src
grpc_dir=staging/src/arcoris.dev/healthgrpc/src
test_dir=staging/src/arcoris.dev/healthtest/src

for forbidden in \
	'"arcoris.dev/healtheval"' \
	'"arcoris.dev/healthprobe"' \
	'"arcoris.dev/healthhttp"' \
	'"arcoris.dev/healthgrpc"' \
	'"arcoris.dev/healthtest"' \
	'"arcoris.dev/chrono' \
	'"arcoris.dev/snapshot' \
	'"google.golang.org/grpc' \
	'"google.golang.org/protobuf'
do
	if has_go_import "$health_dir" "$forbidden"; then
		fail "health must not import $forbidden"
	fi
done

for forbidden in '"arcoris.dev/healthprobe"' '"arcoris.dev/healthhttp"' '"arcoris.dev/healthgrpc"' '"google.golang.org/grpc' '"google.golang.org/protobuf'; do
	if has_go_import "$healtheval_dir" "$forbidden"; then
		fail "healtheval must not import $forbidden"
	fi
done

for forbidden in '"arcoris.dev/healtheval"' '"arcoris.dev/healthhttp"' '"arcoris.dev/healthgrpc"' '"google.golang.org/grpc' '"google.golang.org/protobuf'; do
	if has_go_import "$probe_dir" "$forbidden"; then
		fail "healthprobe must not import $forbidden"
	fi
done

for forbidden in '"arcoris.dev/healtheval"' '"arcoris.dev/healthprobe"' '"arcoris.dev/healthgrpc"' '"arcoris.dev/chrono' '"arcoris.dev/snapshot' '"google.golang.org/grpc' '"google.golang.org/protobuf'; do
	if has_go_import "$http_dir" "$forbidden"; then
		fail "healthhttp must not import $forbidden"
	fi
done

for forbidden in '"arcoris.dev/healtheval"' '"arcoris.dev/healthprobe"' '"arcoris.dev/healthhttp"'; do
	if has_go_import "$grpc_dir" "$forbidden"; then
		fail "healthgrpc must not import $forbidden"
	fi
done

for forbidden in '"arcoris.dev/healtheval"' '"arcoris.dev/healthprobe"' '"arcoris.dev/healthhttp"' '"arcoris.dev/healthgrpc"' '"arcoris.dev/chrono' '"arcoris.dev/snapshot' '"google.golang.org/grpc' '"google.golang.org/protobuf'; do
	if has_go_import "$test_dir" "$forbidden"; then
		fail "healthtest must not import $forbidden"
	fi
done

for mod in "$probe_dir/go.mod" "$http_dir/go.mod" "$grpc_dir/go.mod" "$test_dir/go.mod"; do
	if has_mod_ref "$mod" 'arcoris.dev/healtheval'; then
		fail "$mod must not require or replace arcoris.dev/healtheval"
	fi
done

for mod in "$health_dir/go.mod" "$healtheval_dir/go.mod" "$probe_dir/go.mod" "$http_dir/go.mod" "$test_dir/go.mod"; do
	if has_mod_ref "$mod" 'google.golang.org/grpc'; then
		fail "only healthgrpc may require google.golang.org/grpc"
	fi
done

if ! has_mod_ref "$grpc_dir/go.mod" 'google.golang.org/grpc'; then
	fail "healthgrpc must require google.golang.org/grpc"
fi
