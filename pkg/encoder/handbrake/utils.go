package handbrake

import "strconv"

func inSliceString(items []string, find string) bool {
	for _, item := range items {
		if item == find {
			return true
		}
	}

	return false
}

func intToSlice(val int) (out []string) {

	for i := 0; i < val; i++ {
		out = append(out, strconv.Itoa(i))
	}

	return
}
