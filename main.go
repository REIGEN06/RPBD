package main

import (
	"console_game/functions"
	"fmt"
)

func main() {
	creature := functions.NewPlayer()
	fmt.Printf("Player created: Hole %v Health  %v Respect  %v Weight  %v\n",
		creature.Hole, creature.Health, creature.Respect, creature.Weight)

	creature.Day()
	fmt.Println(creature)
}
