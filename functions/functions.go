package functions

import (
	"fmt"
)

func Hello() {
	fmt.Println("Hello, World!")
}

type creature struct {
	Hole    int16
	Health  int16
	Respect int16
	Weight  int16
}

func NewPlayer() creature {
	return creature{
		Hole:    8,
		Health:  100,
		Respect: 10,
		Weight:  20,
	}
}
func (c *creature) Day() {
	var input int8
	fmt.Print("Choose the action:\n" +
		"1. Dig the hole\n" +
		"2. Eat a grass\n" +
		"3. Fight\n" +
		"4. Sleep\n" +
		"Your choice: ")
	fmt.Scan(&input)
	switch input {
	case 1:
		c.dig()
		// case 2:
		// 	c.eat()
		// case 3:
		// 	c.fight()
		// case 4:
		// 	c.sleep()
	}
}

func (c *creature) Night() {
	c.Hole -= 2
	c.Health += 20
	c.Respect -= 2
	c.Weight -= 5
}

func (c *creature) dig() {
	var input int8
	fmt.Print("Choose:\n" +
		"1. High\n" +
		"2. Lazily\n" +
		"Your choice: ")
	fmt.Scan(&input)
	switch input {
	case 1:
		c.Hole += 5
		c.Health -= 30
	case 2:
		c.Hole += 2
		c.Health -= 10
	}
}
