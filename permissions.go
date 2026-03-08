package dispatch

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// PermissionNames maps Discord permission bit values to human readable names.
var PermissionNames = map[int64]string{
	// Channel
	discordgo.PermissionViewChannel:        "View Channel",
	discordgo.PermissionSendMessages:       "Send Messages",
	discordgo.PermissionSendTTSMessages:    "Send TTS Messages",
	discordgo.PermissionManageMessages:     "Manage Messages",
	discordgo.PermissionEmbedLinks:         "Embed Links",
	discordgo.PermissionAttachFiles:        "Attach Files",
	discordgo.PermissionReadMessageHistory: "Read Message History",
	discordgo.PermissionMentionEveryone:    "Mention Everyone",
	discordgo.PermissionAddReactions:       "Add Reactions",

	// Moderation / Management
	discordgo.PermissionKickMembers:     "Kick Members",
	discordgo.PermissionBanMembers:      "Ban Members",
	discordgo.PermissionModerateMembers: "Timeout Members",
	discordgo.PermissionAdministrator:   "Administrator",
	discordgo.PermissionManageNicknames: "Manage Nicknames",
	discordgo.PermissionManageRoles:     "Manage Roles",
	discordgo.PermissionManageChannels:  "Manage Channels",
	discordgo.PermissionManageGuild:     "Manage Server",
	discordgo.PermissionViewAuditLogs:   "View Audit Log",

	// General
	discordgo.PermissionCreateInstantInvite: "Create Invite",
}

func permissionName(perm int64) string {
	if name, ok := PermissionNames[perm]; ok {
		return name
	}

	return "Unknown Permission"
}


// hasRequiredPermissions runs the permission validation for the user & bot in the current channel
func hasRequiredPermissions(cmd *Command, data *Data) error {
	err := hasUserPermissions(cmd, data)
	if err != nil {
		return err
	}

	err = hasBotPermissions(cmd, data)
	if err != nil {
		return err
	}

	return nil
}

// hasUserPermissions validates that the user has one of any permissions required
func hasUserPermissions(cmd *Command, data *Data) error {
	if len(cmd.RequiredUserPerms) < 1 {
		return nil
	}

	perms, err := data.Session.State.UserChannelPermissions(data.Author.ID, data.Channel.ID)
	if err != nil {
		return fmt.Errorf("Failed retrieving the users perms in the current channel.")
	}

	permsMet := false
	for _, perm := range cmd.RequiredUserPerms {
		if perms&perm != 0 {
			permsMet = true
			break
		}
	}

	if !permsMet {
		humanisedPerms := humanisePerms(cmd.RequiredUserPerms)
		return fmt.Errorf("You do not have any of the required permissions to run this command.\nThis command required one of the following permissions: %s", humanisedPerms)
	}

	return nil
}

// hasBotPermissions validates that the bot has one of any permissions required
func hasBotPermissions(cmd *Command, data *Data) error {
	if len(cmd.RequiredBotPerms) < 1 {
		return nil
	}

	perms, err := data.Session.State.UserChannelPermissions(data.Bot.ID, data.Channel.ID)
	if err != nil {
		return fmt.Errorf("Failed retrieving the bots perms in the current channel.")
	}

	permsMet := false
	for _, perm := range cmd.RequiredBotPerms {
		if perms&perm != 0 {
			permsMet = true
			break
		}
	}

	if !permsMet {
		humanisedPerms := humanisePerms(cmd.RequiredBotPerms)
		return fmt.Errorf("The bot does not have any of the required permissions to run this command.\nThis command required one of the following permissions: %s", humanisedPerms)
	}

	return nil
}

func humanisePerms(perms []int64) string {
	humanisedPerms := make([]string, 0, len(perms))
	for _, perm := range perms {
		humanisedPerms = append(humanisedPerms, fmt.Sprintf("`%s`", permissionName(perm)))
	}

	return strings.Join(humanisedPerms, " or ")
}
