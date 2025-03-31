package ptr

func Get[T any](v *T, def ...T) T {
	if v == nil {
		if len(def) > 0 {
			return def[0]
		}
		var zero T
		return zero
	}
	return *v
}

func To[T any](v T) *T {
	return &v
}
