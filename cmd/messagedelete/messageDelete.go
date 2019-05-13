package messagedelete

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// Entrypoint is the entrypoint into all the delete
// message functions
func Entrypoint(s *discordgo.Session, m *discordgo.MessageCreate) (err error) {

	content := m.Content[2:len(m.Content)]
	content = content[len("delete "):len(content)]
	fmt.Println(content)
	s.ChannelMessageSend(m.ChannelID, "Went through the entrypoint I guess")

	return nil
}
