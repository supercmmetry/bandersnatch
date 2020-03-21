package main

import (
	"bandersnatch/pkg/game"
	"bufio"
	"fmt"
	"os"
)

func main() {
	fmt.Println("Bandersnatch: A Dynamically Randomized State Automaton (a.k.a DYRASTAT)")
	nexus := &game.Nexus{}
	if err := nexus.LoadFromFile("sample.json"); err != nil {
		fmt.Println(err)
	}

	p := &game.Player{Id: 1729}
	nexus.Start(p)
	r := bufio.NewReader(os.Stdin)
	for true {
		fmt.Println(p.CurrentNode.Data.Question)
		fmt.Println("1.", p.CurrentNode.Data.LeftOption)
		fmt.Println("2.", p.CurrentNode.Data.RightOption)
		opt, _ := r.ReadString('\n')
		opt = opt[:1]
		if opt == "1" {
			nexus.Traverse(p, game.OptionLeft)
		} else if opt == "2" {
			nexus.Traverse(p, game.OptionRight)
		} else {
			break
		}

		if artifact := nexus.CheckForArtifact(p); artifact != nil {
			fmt.Println("You found a bandersnatch-artifact.")
			fmt.Println("Artifact description: ", artifact.Description)
		}
	}
}
