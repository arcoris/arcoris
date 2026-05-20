module arcoris.dev/resilience

go 1.25.0

toolchain go1.25.9

require (
	arcoris.dev/admission v0.0.0
	arcoris.dev/capacity v0.0.0
	arcoris.dev/chrono v0.0.0
	arcoris.dev/snapshot v0.0.0
)

replace arcoris.dev/admission => ../admission
replace arcoris.dev/capacity => ../capacity
replace arcoris.dev/chrono => ../chrono
replace arcoris.dev/snapshot => ../snapshot
