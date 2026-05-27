#!/bin/sh

set -eu

repo_root=$(
	cd -- "$(dirname -- "$0")/.." >/dev/null 2>&1
	pwd
)

cd "$repo_root"

status=0

fail() {
	printf '%s\n' "$*" >&2
	status=1
}

manifests=$(find staging/src/arcoris.dev -mindepth 2 -maxdepth 2 -name publishing.yaml | sort)

for manifest in $manifests; do
	module_root=$(dirname "$manifest")
	module=$(sed -n 's/^module: //p' "$manifest")
	src_root="$module_root/src"
	go_mod="$src_root/go.mod"

	if [ ! -f "$go_mod" ]; then
		fail "missing $go_mod"
		continue
	fi

	module_path=$(sed -n 's/^module //p' "$go_mod" | head -n 1)
	if [ "$module_path" != "$module" ]; then
		fail "$go_mod module path $module_path does not match $module"
	fi

	if [ -f "$module_root/go.mod" ]; then
		fail "unexpected parent-level go.mod: $module_root/go.mod"
	fi

	if ! grep -Fq "./$module_root/src" go.work; then
		fail "go.work is missing ./$module_root/src"
	fi

	if rg -n "\"$module/src" "$src_root" >/dev/null 2>&1; then
		fail "forbidden import path contains /src under $src_root"
	fi
done

exit "$status"
