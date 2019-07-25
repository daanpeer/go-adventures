package main

import "fmt"

func pull(slice []string, subject string) (string, error) {
	for _, element := range slice {
		fmt.Println(element, subject)
		if element == subject {
			println("equal!")
			return element, nil
		}
	}
	return "", nil
}
