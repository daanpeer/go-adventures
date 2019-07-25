package main

import (
	"time"
)

type Todo struct {
	ID          int
	Description string
	Date        time.Time
}

type Todos struct {
	todos []Todo
}

func (t *Todos) removeTodo(ID int) {
	var indexToRemove int

	for index, element := range t.todos {
		if element.ID == ID {
			indexToRemove = index
			break
		}
	}

	t.todos = t.todos[:indexToRemove+copy(t.todos[indexToRemove:], t.todos[indexToRemove+1:])]
}

func (t *Todos) hasTodo(ID int) bool {
	for _, element := range t.todos {
		if ID == element.ID {
			return true
		}
	}
	return false
}

func (t *Todos) addTodo(todo Todo) {
	t.todos = append(t.todos, todo)
}
