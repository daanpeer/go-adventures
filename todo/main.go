package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

func main() {

	todos := Todos{}

	http.HandleFunc("/todos/", func(writer http.ResponseWriter, request *http.Request) {
		if request.Method == http.MethodDelete {
			id, error := strconv.Atoi(request.URL.Path[len("/todos/"):])
			if error != nil {
				writer.WriteHeader(403)
				writer.Write([]byte("Unprocessable entity"))
				return
			}

			if !todos.hasTodo(id) {
				writer.WriteHeader(404)
				writer.Write([]byte("Item not found"))
				return
			}

			todos.removeTodo(id)
			return
		}

		if request.Method == http.MethodPost {
			body, err := ioutil.ReadAll(request.Body)

			newTodo := Todo{}

			err = json.Unmarshal(body, &newTodo)

			if err != nil {
				writer.WriteHeader(403)
				writer.Write([]byte("Unprocessable entity"))
				return
			}

			newID := len(todos.todos)

			todos.addTodo(Todo{ID: newID + 1, Description: newTodo.Description, Date: newTodo.Date})
			responseData, _ := json.Marshal(todos.todos)
			writer.Write(responseData)
			return
		}

		responseData, _ := json.Marshal(todos.todos)
		writer.Write(responseData)
	})

	fmt.Print("starting server on port 8081")

	error := http.ListenAndServe(":8081", nil)

	fmt.Println(error)
}
