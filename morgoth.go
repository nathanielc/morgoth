package main

import (
	"flag"
	"fmt"
	"github.com/golang/glog"
	"github.com/nvcook42/morgoth/app"
	"github.com/nvcook42/morgoth/config"
	"os"
)

var configPath = flag.String("config", "morgoth.yaml", "Path to morgoth config")

func main() {
	defer glog.Flush()
	flag.Parse()
	config, err := config.LoadFromFile(*configPath)
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
