module arcoris.dev/healthprobe

go 1.25.0

toolchain go1.25.9

require (
	arcoris.dev/chrono v0.0.0
	arcoris.dev/health v0.0.0
	arcoris.dev/healthtest v0.0.0
	arcoris.dev/snapshot v0.0.0
)

require (
	arcoris.dev/atomicx v0.0.0 // indirect
	arcoris.dev/healtheval v0.0.0 // indirect
)

replace arcoris.dev/atomicx => ../../atomicx/src

replace arcoris.dev/chrono => ../../chrono/src

replace arcoris.dev/health => ../../health/src

replace arcoris.dev/healtheval => ../../healtheval/src

replace arcoris.dev/healthtest => ../../healthtest/src

replace arcoris.dev/snapshot => ../../snapshot/src

replace arcoris.dev/testutil => ../../testutil/src
