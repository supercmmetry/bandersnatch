package main

import (
	bs "bandersnatch/core"
	"bufio"
	"fmt"
	"os"
)

func main() {
	fmt.Println("Bandersnatch (Dynamically Randomized State Automaton)")
	nexus := &bs.Nexus{}
	if err := nexus.LoadFromFile("sample.json"); err != nil {
		fmt.Println(err)
	}

	p := &bs.Player{Id:1729}
	nexus.Start(p)
	r := bufio.NewReader(os.Stdin)
	for true {
		fmt.Println(p.CurrentNode.Data.Question)
		fmt.Println("1.", p.CurrentNode.Data.LeftOption)
		fmt.Println("2.", p.CurrentNode.Data.RightOption)
		opt, _ := r.ReadString('\n')
		opt = opt[:1]
		if opt == "1" {
			nexus.Traverse(p, bs.OptionLeft)
		} else if opt == "2" {
			nexus.Traverse(p, bs.OptionRight)
		} else {
			break
		}
	}
}
