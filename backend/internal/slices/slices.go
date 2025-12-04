package slices

func Delete[T any](slice []*T, needle *T) {
	for i := range slice {
		if needle == slice[i] {
			slice = append(slice[:i], slice[i+1:]...)
			break
		}
	}
}
