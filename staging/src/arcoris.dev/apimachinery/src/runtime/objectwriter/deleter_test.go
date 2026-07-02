package objectwriter

import "arcoris.dev/apimachinery/api/objectlifecycle"

var _ Deleter = (*objectlifecycle.Executor)(nil)
