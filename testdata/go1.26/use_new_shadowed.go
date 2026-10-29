package pkg

// Nothing is reported in this file: the function new declared below hides the built-in one,
// so new(value) would call it instead of allocating.

func pointerToSize(i int) {
	size := float64(i)
	snapshot.Size = &size
}

func new() *builderType {
	return &builderType{}
}
