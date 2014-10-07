package main

import (
	"fmt"
	"github.com/nvcook42/morgoth/app"
	"github.com/nvcook42/morgoth/config"
	"os"
)

func main() {
	config, err := config.LoadFromFile(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	app, err := app.New(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error during startup: %v\n", err)
		os.Exit(2)
	}
	err = app.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running application: %v\n", err)
		os.Exit(3)
	}
}
