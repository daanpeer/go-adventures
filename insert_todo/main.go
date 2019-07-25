package main

import (
	"bytes"
	"fmt"
	"net/http"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	amount := 100
	fmt.Println("Inserting", amount, "records")
	data := []byte(`{"decription":"test123"}`)
	wg.Add(amount)
	for index := 0; index < amount; index++ {
		go func() {
			_, err := http.Post("http://localhost:8081/todos/", "application/json", bytes.NewBuffer(data))

			// fmt.Println(resp)

			if err != nil {
				fmt.Println(err)
			}

			defer wg.Done()
		}()
	}

	fmt.Println("Done creating requests!")
	wg.Wait()

	fmt.Println("Requests sent!")
}
