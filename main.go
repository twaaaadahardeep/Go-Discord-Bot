package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var (
	Token string
)

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {
	discord, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("Failed here")
		log.Fatal(err)
	}

	discord.AddHandler(messageCreate)
	discord.AddHandler(presenceUpdateHandler)

	discord.Identify.Intents = discordgo.IntentsGuildMessages + discordgo.IntentsGuildPresences

	e := discord.Open()
	if e != nil {
		log.Fatal(e)
	}

	fmt.Println("Bot is running now. Press Ctrl + C to exit")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	discord.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "ping" {
		s.ChannelMessageSend(m.ChannelID, "pong")
		fmt.Println(m.ChannelID)
	}

	if m.Content == "pong" {
		s.ChannelMessageSend(m.ChannelID, "ping")
	}
}

func presenceUpdateHandler(s *discordgo.Session, p *discordgo.PresenceUpdate) {
	if p.Presence.Status == discordgo.StatusOnline {
		ch, err := s.UserChannelCreate(p.User.ID)
		if err != nil {
			log.Fatal("Error creating User Channel")
		}

		s.ChannelMessageSend(ch.ID, "Hi Welcome.")
	}
}
