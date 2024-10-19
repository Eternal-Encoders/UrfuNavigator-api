package utils

func MapToArray[Q comparable, T any](args map[Q]T) []T {
	arr := make([]T, len(args))
	i := 0
	for _, value := range args {
		arr[i] = value
		i++
	}
	return arr
}
