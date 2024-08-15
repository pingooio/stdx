package slicesx

func Map[S ~[]E, E any](input S, fn func(elem E, index int) E) S {
	ret := make([]E, len(input))

	for i, elem := range input {
		ret[i] = fn(elem, i)
	}

	return ret
}
