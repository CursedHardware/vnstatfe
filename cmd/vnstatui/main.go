package main

import (
	"flag"
	"github.com/CursedHardware/vnstatfe/vnstat"
	"log"
	"net/http"
)

var arguments []string
var host string

func init() {
	flag.StringVar(&host, "host", ":7991", "Host")
	configFile := flag.String("config", "", "path to config file")
	iface := flag.String("iface", "", "network interface to use")
	flag.Parse()
	if *configFile != "" {
		arguments = append(arguments, "--config", *configFile)
	}
	if *iface != "" {
		arguments = append(arguments, "--iface", *iface)
	}
}

func main() {
	handler, err := vnstat.NewHandler(arguments)
	if err != nil {
		return
	}
	log.Printf("Running on %s\n", host)
	http.Handle("/", handler)
	if err = http.ListenAndServe(host, nil); err != nil {
		log.Fatalln(err)
	}
}
