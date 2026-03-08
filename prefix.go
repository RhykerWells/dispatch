package dispatch

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

// checkMessagePrefix checks a given content for the prefix or bot mention of the guild
func (c *CommandHandler) checkMessagePrefix(session *discordgo.Session, event *discordgo.MessageCreate) (string, bool) {
	if prefix, ok := findBasicPrefix(c.Prefix(event.GuildID), event.Content); ok {
		return prefix, ok
	}

	if prefix, ok := findMentionPrefix(session.State.User.ID, event.Content); ok {
		return prefix, ok
	}

	return "", false
}

// findBasicPrefix finds a text based prefix such as "-" or "~"
func findBasicPrefix(prefix, message string) (string, bool) {
	if !strings.HasPrefix(message, prefix) {
		return "", false
	}

	return prefix, true
}

// findMentionPrefix finds a bot mention prefix such as @Discord
func findMentionPrefix(botID string, message string) (string, bool) {
	mentionPrefixes := []string{"<@" + botID + ">", "<@!" + botID + ">"}
	for _, mentionPrefix := range mentionPrefixes {
		if strings.HasPrefix(message, mentionPrefix) {
			return mentionPrefix, true
		}
	}

	return "", false
}