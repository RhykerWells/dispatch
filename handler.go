package dispatch

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

// CommandHandler contains the prefix, the full instances of a command and a string map to retireve them
type CommandHandler struct {
	Prefix       func(string) string
	cmdInstances []*Command
	cmdMap       map[string]*Command
}

// NewCommandHandler creates a new command handler
func NewCommandHandler() *CommandHandler {
	handler := &CommandHandler{
		Prefix:       func(string) string { return "~" },
		cmdInstances: make([]*Command, 0),
		cmdMap:       make(map[string]*Command),
	}

	return handler
}

// SetPrefixFunc allows you to set a custom function to determine the prefix for a guild
func (c *CommandHandler) SetPrefixFunc(prefixFunc func(string) string) {
	c.Prefix = prefixFunc
}

// RegisterCommands adds each command to the command handler
func (c *CommandHandler) RegisterCommands(cmds ...*Command) {
	for _, cmd := range cmds {
		c.cmdInstances = append(c.cmdInstances, cmd)

		cleanAliases := make([]string, 0, len(cmd.Aliases))
		for _, alias := range cmd.Aliases {
			if trimmed := strings.TrimSpace(alias); trimmed != "" {
				cleanAliases = append(cleanAliases, trimmed)
			}
		}

		if len(cleanAliases) > 3 {
			aliasOver := len(cleanAliases) - 3
			cleanAliases = cleanAliases[:3]
			logrus.Warnf("%s has %d too many non-empty aliases. Automatically removed the last %d.", cmd.Command, aliasOver, aliasOver)
		}

		// Register main command
		c.cmdMap[strings.ToLower(cmd.Command)] = cmd

		// Register aliases
		for _, alias := range cleanAliases {
			c.cmdMap[strings.ToLower(alias)] = cmd
		}
	}
}

// RegisteredCommands returns an array of each RegisteredCommand
func (c *CommandHandler) RegisteredCommands() map[string]RegisteredCommand {
	cmdMap := make(map[string]RegisteredCommand, len(c.cmdInstances))

	for _, cmd := range c.cmdInstances {
		cmdMap[cmd.Command] = RegisteredCommand{Command: cmd}
	}

	return cmdMap
}

// Handles all message create events to the bot, to pass them to child functions
func (c *CommandHandler) HandleMessageCreate(s *discordgo.Session, e *discordgo.MessageCreate) {
	if e.Author.ID == s.State.User.ID || e.Author.Bot {
		return
	}

	prefix, ok := c.checkMessagePrefix(s, e)
	if !ok {
		return
	}

	args := strings.Fields(e.Content[len(prefix):])
	if len(args) < 1 {
		return
	}

	command := strings.ToLower(args[0])
	cmd, ok := c.cmdMap[command]
	if !ok {
		return
	}

	guild, err := s.State.Guild(e.GuildID)
	if err != nil {
		guild, _ = s.Guild(e.GuildID)
	}

	channel, err := s.State.Channel(e.ChannelID)
	if err != nil {
		channel, _ = s.Channel(e.ChannelID)
	}

	data := &Data{
		Session:    s,
		Guild:      guild,
		Channel:    channel,
		Message:    e.Message,
		Author:     e.Author,
		ParsedArgs: nil,
		Handler:    c,
	}

	go runCommand(cmd, data, args[1:])
}

// runCommand logs the command called by the bot, ensures the correct number of args is present, parses the args, then runs the command
func runCommand(cmd *Command, data *Data, tokens []string) {
	logrus.WithFields(logrus.Fields{
		"Guild":           data.Guild.ID,
		"Command":         cmd.Command,
		"Triggering user": data.Author.ID},
	).Infoln("Executed command")

	err := hasRequiredPermissions(cmd, data)
	if err != nil {
		data.Session.ChannelMessageSend(data.Channel.ID, err.Error())
		return
	}

	err = parseArgs(cmd, data, tokens)
	if err != nil {
		data.Session.ChannelMessageSend(data.Channel.ID, err.Error())
		return
	}

	err = cmd.Run(data)
	if err != nil {
		data.Session.ChannelMessageSend(data.Channel.ID, err.Error())
	}
}
