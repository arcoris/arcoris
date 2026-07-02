package objectwriter

import "arcoris.dev/apimachinery/api/objectlifecycle"

var _ MetadataPatcher = (*objectlifecycle.Executor)(nil)
