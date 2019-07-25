package main

func fill(slice []string, content string) []string {
	for index := 0; index < len(slice); index++ {
		slice[index] = content
	}
	return slice
}
