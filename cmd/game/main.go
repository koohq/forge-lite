package main

import (
	"os"

	"github.com/koohq/forge-lite/internal/engine"
)

func main() {
	game := engine.NewGame(os.Stdin)
	game.Start()
}
