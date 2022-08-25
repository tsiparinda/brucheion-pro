package repo

//removeDuplicatesUnordered takes a string slice and returns it without any duplicates.
func RemoveDuplicatesUnordered(elements []string) []string {
	encountered := map[string]bool{}

	// Create a map of all unique elements.
	for v := range elements {
		encountered[elements[v]] = true
	}

	// Place all keys from the map into a slice.
	result := []string{}
	for key := range encountered {
		result = append(result, key)
	}
	return result
}

//contains returns true if the 'needle' string is found in the 'heystack' string slice
func Contains(heystack []string, needle string) bool {
	for _, straw := range heystack {
		if straw == needle {
			return true
		}
	}
	return false
}
