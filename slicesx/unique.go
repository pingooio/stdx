package slicesx

func Unique[T comparable](input []T) []T {
	inResult := make(map[T]bool, len(input))
	var result []T

	for _, elem := range input {
		if _, ok := inResult[elem]; !ok {
			inResult[elem] = true
			result = append(result, elem)
		}
	}

	return result
}
