package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHas(t *testing.T) {
	items := [...]string{"a", "b", "c"}
	item, _ := pull(items[:], "a")
	assert.Equal(t, item, "a")
}
