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
	"github.com/cull-methi/malachite/cmd/messagedelete"
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
	configFilePath := os.Getenv("CONFIG_FILE")
	creds := Creds{}
	err := fillCredsFromConfigFile(configFilePath, &creds)
	if err != nil {
		log.Fatal(err)
	}

	dg, err := discordgo.New("Bot " + creds.BotToken)
	if err != nil {
		log.Fatal(err)
	}

	dg.AddHandler(messageCreateHandler)

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

func fillCredsFromConfigFile(configFilePath string, creds *Creds) (err error) {
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

func messageCreateHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

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
		s.ChannelMessageSend(m.ChannelID, "Yeah, yeah, I'm here bruh")
	}

	if content == "test" {
		s.ChannelMessageSend(m.ChannelID, "Here's your test, fuckwad")
	}

	err := error{}
	if strings.HasPrefix(content, "delete") {
		err = messagedelete.Entrypoint(s, m)
	}

}
