package main

import (
	"os"

	"forge-lite/internal/engine"
)

func main() {
	game := engine.NewGame(os.Stdin)
	game.Start()
}
