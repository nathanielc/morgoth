package main

import (
	"flag"
	"github.com/nathanielc/morgoth"
	"github.com/nathanielc/morgoth/Godeps/_workspace/src/github.com/golang/glog"
	"os"
)

var configPath = flag.String("config", "morgoth.yaml", "Path to morgoth config")

func main() {
	defer glog.Flush()
	flag.Parse()
	config, err := morgoth.LoadFromFile(*configPath)
	if err != nil {
		glog.Errorf("Error loading config: %v\n", err)
		os.Exit(2)
	}

	app := morgoth.NewApp(config)
	err = app.Run()
	if err != nil {
		glog.Errorf("Error running application: %v\n", err)
		os.Exit(3)
	}
}
