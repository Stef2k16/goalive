package main

import (
	"fmt"
	"galive/notification"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	token := "ODIyODQ4MjQyNDYyMzU5NTky.YFYPJA.1KM0nt7xmRpHt8aiqYyi1bQJnQQ"
	channelID := "822848790955294723"
	var c notification.Client
	c, err := notification.NewDiscordClient(token, channelID)
	if err != nil {
		log.Fatal("Error creating session: ", err)
		return
	}

	fmt.Println("Notification Bot is now running. Press CTRL-C to exit.")

	err = c.SendNotification("I'm alive!")
	fmt.Println(err)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	if err := c.Close(); err != nil {
		log.Fatal("Error closing session: ", err)
	}
}
