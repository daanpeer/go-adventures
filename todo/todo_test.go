package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTodos_removeTodo(test *testing.T) {
	var t Todos

	time := time.Now()

	items := []Todo{
		{ID: 2, Description: "test", Date: time},
		{ID: 12, Description: "test", Date: time},
		{ID: 90, Description: "test", Date: time},
	}

	for _, element := range items {
		t.addTodo(element)
	}

	results := []Todo{
		{ID: 2, Description: "test", Date: time},
		{ID: 90, Description: "test", Date: time},
	}

	t.removeTodo(12)

	assert.Equal(test, t.todos, results)
}

func TestTodos_addTodo(test *testing.T) {
	var t Todos

	todoToAdd := Todo{ID: 1, Description: "test", Date: time.Now()}

	todos := []Todo{todoToAdd}

	t.addTodo(todoToAdd)

	assert.Equal(test, t.todos, todos)
}
