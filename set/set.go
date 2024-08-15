package set

type Set[T comparable] map[T]struct{}

func New[T comparable](fromList []T) Set[T] {
	ret := Set[T](make(map[T]struct{}))

	for _, item := range fromList {
		ret[item] = struct{}{}
	}

	return ret
}

func (set Set[T]) Contains(item T) bool {
	_, contains := set[item]
	return contains
}

func (set Set[T]) ToSlice() []T {
	ret := make([]T, len(set))
	i := 0

	for elem := range set {
		ret[i] = elem
		i += 1
	}

	return ret
}

func (set Set[T]) Add(element T) {
	set[element] = struct{}{}
}
