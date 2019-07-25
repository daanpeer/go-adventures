package main

import (
	"fmt"
)

type Person struct {
	Name string
}

type Saiyan struct {
	*Person
	Power int
}

func (s *Saiyan) Super() {
	s.Power += 100090
}

func (p *Person) Introduce() {
	fmt.Println("Hello person!", p.Name)
}

func (s *Saiyan) Introduce() {
	fmt.Println("Hello saiyan!", s.Name)
}

func main() {
	saiyan := &Saiyan{
		Person: &Person{Name: "Goku"},
		Power:  9000,
	}

	saiyan.Super()
	saiyan.Introduce()
	saiyan.Person.Introduce()

	fmt.Println("blala", saiyan)
}
