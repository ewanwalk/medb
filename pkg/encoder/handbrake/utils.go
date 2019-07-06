package handbrake

import (
	"strconv"
	"strings"
)

func inSliceString(items []string, find string) bool {
	for _, item := range items {
		if item == find {
			return true
		}
	}

	return false
}

func inSliceInt(items []int, find int) bool {
	for _, item := range items {
		if item == find {
			return true
		}
	}

	return false
}

func intToSlice(val int) (out []string) {

	for i := 0; i < val; i++ {
		out = append(out, strconv.Itoa(i+1))
	}

	return
}

func repeatInt(val, repeat int) string {

	out := make([]string, 0, repeat)
	strVal := strconv.Itoa(val)
	for i := 0; i < repeat; i++ {
		out = append(out, strVal)
	}

	return strings.Join(out, ",")
}
