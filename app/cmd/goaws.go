package main

import (
	"flag"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/jeff-blaisdell/goaws/app/conf"
	"github.com/jeff-blaisdell/goaws/app/router"
)

func main() {
	var filename string
	var debug bool
	flag.StringVar(&filename, "config", "", "config file location + name")
	flag.BoolVar(&debug, "debug", false, "debug log level (default Warning)")
	flag.Parse()
	log.Warnf("file: %s\n", filename)

	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	if debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.WarnLevel)
	}

	env := "Local"
	if flag.NArg() > 0 {
		log.Warn(flag.Args())
		env = flag.Arg(0)
	}

	log.Warnf("file: %s, Env: %s\n", filename, env)
	portNumbers := conf.LoadYamlConfig(filename, env)

	r := router.New()

	if len(portNumbers) == 1 {
		log.Warnf("GoAws listening on: 0.0.0.0:%s", portNumbers[0])
		err := http.ListenAndServe("0.0.0.0:"+portNumbers[0], r)
		log.Fatal(err)
	} else if len(portNumbers) == 2 {
		go func() {
			log.Warnf("GoAws listening on: 0.0.0.0:%s", portNumbers[0])
			err := http.ListenAndServe("0.0.0.0:"+portNumbers[0], r)
			log.Fatal(err)
		}()
		log.Warnf("GoAws listening on: 0.0.0.0:%s", portNumbers[1])
		err := http.ListenAndServe("0.0.0.0:"+portNumbers[1], r)
		log.Fatal(err)
	} else {
		log.Fatal("Not enough or too many ports defined to start GoAws.")
	}
}
