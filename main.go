package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
)

var (
	infoLogger    *log.Logger
	warningLogger *log.Logger
	errorLogger   *log.Logger
	cnf *config
)

func init() {
	infoLogger = log.New(os.Stderr, "INFO ", log.LstdFlags)
	warningLogger = log.New(os.Stderr, "WARNING ", log.LstdFlags)
	errorLogger = log.New(os.Stderr, "ERROR ", log.LstdFlags)
}

func main()  {
	configFile := flag.String("config", "config.yml", "config file")
	flag.Parse()

	flag.Parse()

	configString := ""
	if *configFile != "" {
		configBytes, err := ioutil.ReadFile(*configFile)
		// error means that the file might not exists, in which case we'll just make it later
		if err == nil {
			configString = string(configBytes)
		}
	}

	var err error
	cnf, err = getConfig(configString, configFile)
	if err != nil {
		errorLogger.Fatalf("Failed to get config: %v", err)
	}

	// setup modules
	setupWebApi(cnf)
}
