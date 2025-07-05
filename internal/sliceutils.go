package internal

func transform[I, O any](ii []I, f func(i I) (O, error)) ([]O, error) {
	oo := make([]O, len(ii))
	for i, item := range ii {
		o, err := f(item)
		if err != nil {
			return nil, err
		}

		oo[i] = o
	}

	return oo, nil
}
