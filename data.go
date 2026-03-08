package dispatch

import (
	"github.com/bwmarrin/discordgo"
)

// Data defines the required data passed to each command
type Data struct {
	Session *discordgo.Session
	Bot     *discordgo.User

	Guild   *discordgo.Guild
	Channel *discordgo.Channel
	Author  *discordgo.User

	Message    *discordgo.Message
	ParsedArgs []*ParsedArg

	Handler *CommandHandler
}
