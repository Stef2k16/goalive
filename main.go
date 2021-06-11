package main

import (
	"flag"
	"fmt"
	"github.com/Stef2k16/goalive/config"
	"github.com/Stef2k16/goalive/monitor"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	configPath := parseFlags()
	conf, err := config.New(configPath)
	if err != nil {
		log.Fatal("Error initializing configuration: ", err)
	}

	m, err := monitor.New(conf)
	if err != nil {
		log.Fatal("Error creating monitor: ", err)
	}
	m.Start()

	fmt.Println("Monitoring and Notification Bot are now running. Press CTRL-C to exit.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	if err := m.NotificationClient.Stop(); err != nil {
		log.Fatal("Error closing session of the notification client: ", err)
	}
}

// parseFlags parses the command line flags and returns its value.
func parseFlags() string {
	configPath := flag.String("config", "", "path to the configuration file")
	flag.Parse()
	return *configPath
}
