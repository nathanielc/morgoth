package main

import (
	"fmt"
	log "github.com/cihub/seelog"
	"github.com/nvcook42/morgoth/app"
	"github.com/nvcook42/morgoth/config"
	"os"
)

func main() {
	defer log.Flush()
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage:\n")
		os.Exit(1)
	}
	config, err := config.LoadFromFile(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(2)
	}

	app := app.New(config)
	err = app.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running application: %v\n", err)
		os.Exit(3)
	}
}
