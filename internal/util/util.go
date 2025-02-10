package util

func FindIndex(arr []string, target string) int {
	for i, v := range arr {
		if v == target {
			return i
		}
	}
	return -1 // Not found
}
