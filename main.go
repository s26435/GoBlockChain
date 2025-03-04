package main

import (
	. "go.mod/blockchain"
)

func main() {
	bc := NewBlockChain()
	defer bc.DB.Close()

	cli := NewCLI(bc)
	cli.Run()
}