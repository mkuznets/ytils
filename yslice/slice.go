package yslice

func Map[E any, R any](slice []E, mapper func(value E) R) []R {
	mapped := make([]R, 0, len(slice))
	for _, el := range slice {
		mapped = append(mapped, mapper(el))
	}
	return mapped
}

func Unique[R comparable](slice []R) []R {
	unique := make([]R, 0)
	visited := map[R]bool{}

	for _, value := range slice {
		if exists := visited[value]; !exists {
			unique = append(unique, value)
			visited[value] = true
		}
	}
	return unique
}

func UniqueMap[E any, R comparable](slice []E, mapper func(value E) R) []R {
	return Unique(Map(slice, mapper))
}

func MapByKey[T any, R comparable](slice []T, key func(value T) R) map[R][]T {
	mapped := make(map[R][]T, len(slice))
	for _, el := range slice {
		mapped[key(el)] = append(mapped[key(el)], el)
	}
	return mapped
}

func Filter[T any](slice []T, predicate func(value T) bool) []T {
	filtered := make([]T, 0, len(slice))
	for _, el := range slice {
		if predicate(el) {
			filtered = append(filtered, el)
		}
	}
	return filtered
}

func Intersect[T comparable](as ...[]T) []T {
	switch len(as) {
	case 0:
		return nil
	case 1:
		return as[0]
	default:
		result := make([]T, 0)
		h := map[T]struct{}{}

		for _, value := range as[0] {
			h[value] = struct{}{}
		}

		for _, b := range as[1:] {
			hb := map[T]struct{}{}

			for _, v := range b {
				if _, ok := h[v]; ok {
					hb[v] = struct{}{}
				}
			}
			h = hb
		}

		for k := range h {
			result = append(result, k)
		}
		return result
	}
}
