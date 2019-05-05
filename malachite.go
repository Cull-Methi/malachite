package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"gopkg.in/yaml.v2"
)

var prefix = "./"

// Creds struct to read in the creds yaml file
type Creds struct {
	AppToken          string
	AppID             string
	BotToken          string
	BotName           string
	BotDiscriminator  string
	BotPermissionsInt int
}

func main() {
	// configFilePath := "/mnt/c/Users/Mitch/Documents/malachite.yaml"
	configFilePath := os.Getenv("CONFIG_FILE")
	creds := Creds{}
	err := fillCreds(configFilePath, &creds)
	if err != nil {
		log.Fatal(err)
	}

	dg, err := discordgo.New("Bot " + creds.BotToken)
	if err != nil {
		log.Fatal(err)
	}

	dg.AddHandler(messageCreate)

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

func fillCreds(configFilePath string, creds *Creds) (err error) {
	file, err := os.Open(configFilePath)
	if err != nil {
		return err
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, &creds)
	if err != nil {
		return err
	}

	return nil

}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == s.State.User.ID {
		return
	}

	content := m.Content
	if !strings.HasPrefix(content, prefix) {
		return
	}

	content = content[2:len(content)]
	fmt.Println(content)
	// If the message is "ping" reply with "Pong!"
	if content == "ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}

	// If the message is "pong" reply with "Ping!"
	if content == "pong" {
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	}

	if content == "test" {
		s.ChannelMessageSend(m.ChannelID, "here's your test fuckwad")
	}
}
