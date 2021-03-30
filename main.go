package main

import (
	"fmt"
	"galive/config"
	"galive/monitor"
	"galive/notification"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	conf, err := config.New("config.yaml")
	if err != nil {
		log.Fatal("Error creating notification session: ", err)
	}

	c, err := notification.GetClient(conf.Notification)
	if err != nil {
		log.Fatal("Error creating notification session: ", err)
	}

	m, err := monitor.New(conf.LogFile, c, conf.URL, time.Duration(conf.PollingInterval)*time.Second)
	if err != nil {
		log.Fatal("Error creating monitor: ", err)
	}
	m.Start()

	fmt.Println("Notification Bot is now running. Press CTRL-C to exit.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	if err := c.Close(); err != nil {
		log.Fatal("Error closing session of the notification client: ", err)
	}
}
