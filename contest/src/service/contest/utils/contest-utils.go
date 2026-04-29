package contestserviceutils

import "slices"

func CheckDuplicateItem(arr []string) bool {
	if len(arr) <= 1 {
		return false
	}
	slices.Sort(arr)
	for i := 1; i < len(arr); i++ {
		if arr[i] == arr[i-1] {
			return true
		}
	}
	return false
}
