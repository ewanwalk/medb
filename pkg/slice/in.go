package slice

func InString(needle string, haystack []string) bool {

	for _, item := range haystack {
		if item == needle {
			return true
		}
	}

	return false
}

func InInt(needle int, haystack []int) bool {

	for _, item := range haystack {
		if item == needle {
			return true
		}
	}

	return false
}
