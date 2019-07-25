package main

import "fmt"

func main() {
	items := [...]string{"a", "b", "c"}

	// for _, num := range items {
	// 	println(num)
	// }
	result, error := pull(items[:], "b")
	if error != nil {
		fmt.Print("error!")
	}
	fmt.Println("Found", result)

	newItems := [99]string{}
	fill(newItems[:], "test")
	fmt.Println(newItems)

	subj := map[string]interface{}{
		"test":  1,
		"test2": "123",
	}

	fmt.Println(has(subj, "dsdf"))
}
