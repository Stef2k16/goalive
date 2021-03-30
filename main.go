package main

import (
	"fmt"
	"galive/monitor"
	"galive/notification"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	token := "ODIyODQ4MjQyNDYyMzU5NTky.YFYPJA.1KM0nt7xmRpHt8aiqYyi1bQJnQQ"
	channelID := "822848790955294723"
	var c notification.Client
	c, err := notification.NewDiscordClient(token, channelID)
	if err != nil {
		log.Fatal("Error creating notification session: ", err)
		return
	}
	urls := []string{"https://httpstat.us/404", "https://httpstat.us/200"}

	m, err := monitor.New("testlog.txt", c, urls, 15*time.Second)
	if err != nil {
		log.Fatal("Error creating monitor: ", err)
		return
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
