# Dispatch
A lightweight Discord bot command framework for Go. Dispatch simplifies building command-driven Discord bots with support for argument parsing, permission checking, prefix configuration, and command aliases.

## Features
- **Easy Command Registration**: Register commands with a simple, declarative API
- **Automatic Message Routing**: Built-in Discord message event handling and command dispatch
- **Flexible Argument Parsing**: Support for multiple argument types (strings, integers, users, members, durations, custom types)
- **Permission Management**: Enforce user and bot permission requirements per command
- **Customizable Prefixes**: Dynamic prefix configuration on a per-guild basis
- **Command Aliases**: Support up to 3 aliases per command for flexible invocation
- **Argument Combinations**: Define flexible argument requirement patterns

## Installation
```bash
go get github.com/RhykerWells/dispatch
```

## Quick Start

### 1. Create a CommandHandler
```go
package main

import (
	"github.com/RhykerWells/dispatch"
	"github.com/bwmarrin/discordgo"
)

func main() {
	// Create a new command handler
	handler := dispatch.NewCommandHandler()
	
	// Optional: Set a custom prefix function (default is "~")
	handler.SetPrefixFunc(func(guildID string) string {
		// Return prefix based on guild ID (e.g., from database)
		return "!"
	})
}
```

### 2. Define Commands
```go
// Create a simple ping command
pingCmd := &dispatch.Command{
	Command:      "ping",
	Category:     dispatch.CommandCategory{Name: "Utility", Description: "Utility commands"},
	Description:  "Responds with pong",
	Aliases:      []string{"p"},
	ArgsRequired: 0,
	Run: func(data *dispatch.Data) error {
		data.Session.ChannelMessageSend(data.Channel.ID, "Pong!")
		return nil
	},
}

// Create a command with arguments
echoCmd := &dispatch.Command{
	Command:      "echo",
	Category:     dispatch.CommandCategory{Name: "Utility"},
	Description:  "Echoes back your message",
	Args:         []*dispatch.Arg{{Name: "text", Type: dispatch.String}},
	ArgsRequired: 1,
	Run: func(data *dispatch.Data) error {
		// Argument values are interface{} and must be type asserted
		text := data.ParsedArgs[0].Value.(string)
		data.Session.ChannelMessageSend(data.Channel.ID, text)
		return nil
	},
}
```

### 3. Register Commands
```go
handler.RegisterCommands(pingCmd, echoCmd)
```

### 4. Attach to Discord Session
```go
session.AddHandler(handler.HandleMessageCreate)
```

## Argument Types
Dispatch provides several built-in argument types.
> **Note:** Argument values are stored as `interface{}` and must be type asserted to their expected Go type.
> Built-in argument types return the following values:

| Type | Go Type returned | Description |
|------|---------|-------------|
| `dispatch.String` | `string` | Plain text argument with optional string options |
| `dispatch.Int` | `int` | Integer argument with optional min/max bounds |
| `dispatch.Int64` | `int64` | Integer argument with optional min/max bounds |
| `dispatch.User` | `*discordgo.User` | Discord user mention/ID |
| `dispatch.Member` | `*discordgo.Member` | Discord member mention/ID |
| `dispatch.Duration` | `time.Duration` | Duration string (e.g., "5m", "1h30m") |

### Example with Multiple Arguments
```go
kickCmd := &dispatch.Command{
	Command:     "kick",
	Description: "Kick a member with optional reason",
	Args: []*dispatch.Arg{
		{Name: "member", Type: dispatch.Member},
		{Name: "reason", Type: dispatch.String},
	},
	ArgsRequired: 1, // Member required, reason optional
	RequiredBotPerms: []int64{discordgo.PermissionKickMembers},
	Run: func(data *dispatch.Data) error {
		// Argument values are interface{} and must be type asserted
		member := data.ParsedArgs[0].Value.(*discordgo.Member)
		reason := ""
		if len(data.ParsedArgs) > 1 {
			reason = data.ParsedArgs[1].Value.(string)
		}
		
		// Kick the user...
		return nil
	},
}
```

## Permission Management
Enforce permissions with `RequiredUserPerms` and `RequiredBotPerms`:

```go
banCmd := &dispatch.Command{
	Command:          "ban",
	Description:      "Ban a member",
	RequiredUserPerms: []int64{discordgo.PermissionBanMembers},
	RequiredBotPerms:  []int64{discordgo.PermissionBanMembers},
	Args:             []*dispatch.Arg{{Name: "user", Type: dispatch.User}},
	ArgsRequired:     1,
	Run: func(data *dispatch.Data) error {
		// Argument values are interface{} and must be type asserted
		user := data.ParsedArgs[0].Value.(*discordgo.User)
		data.Session.GuildBanCreate(data.Guild.ID, user.ID, 0)
		return nil
	},
}
```

## Command Categories
Organize commands logically using categories:

```go
utilityCategory := dispatch.CommandCategory{
	Name:        "Utility",
	Description: "Helpful utility commands",
}

cmd := &dispatch.Command{
	Command:  "help",
	Category: utilityCategory,
	// ... rest of command
}
```

## Data Available in Commands
Each command's `Run` function receives a `*dispatch.Data` object with:

- **Session**: The Discord session
- **Bot**: The bot user information
- **Guild**: The guild (server) where the command was invoked
- **Channel**: The channel where the command was invoked
- **Author**: The user who invoked the command
- **Message**: The original Discord message
- **ParsedArgs**: Parsed and validated arguments (When accessed, must be type asserted)
- **Handler**: Reference to the command handler

## Example Discord Bot
```go
package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/RhykerWells/dispatch"
	"github.com/bwmarrin/discordgo"
)

func main() {
	session, _ := discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))
	handler := dispatch.NewCommandHandler()

	// Create commands
	pingCmd := &dispatch.Command{
		Command:      "ping",
		Description:  "Ping pong",
		Run: func(data *dispatch.Data) error {
			data.Session.ChannelMessageSend(data.Channel.ID, "Pong! 🏓")
			return nil
		},
	}

	handler.RegisterCommands(pingCmd)
	session.AddHandler(handler.HandleMessageCreate)
	session.Open()

	// Wait for interrupt
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	session.Close()
}
```

## Advanced Usage

### Custom Argument Types
Implement the `ArgumentType` interface to create custom argument validators:
The `ParseArg` method should return both the parsed data & if successfully parsed so that the
result stored in `ParsedArgs[i].Value` can be safely type asserted in the command handler.
```go
type CustomArg struct{}

var _ dispatch.ArgumentType = (*CustomArg)(nil)

func (c *CustomArg) ParseArg(arg *dispatch.ParsedArg, data *dispatch.Data) (any, bool) {
	value := arg.Raw 
	// Validation logic
	// Return the value in whatever type you'd like
	return value, true
}

func (c *CustomArg) Help() string {
	return "Custom argument type"
}
```


### Argument Combinations
Define flexible argument patterns using argument combinations:

```go
cmd := &dispatch.Command{
	Command:     "example",
	Args:        []*dispatch.Arg{...},
	ArgumentCombos: [][]int{
		{0},        // Just first argument
		{0, 1},     // First and second
		{0, 1, 2},  // All three
	},
	Run: func(data *dispatch.Data) error {
		return nil
	},
}
```

## Important Notes
- **Prefix**: Default prefix is `~`. Customize it with `SetPrefixFunc()`
- **Aliases**: Limited to 3 non-empty aliases per command
- **Bot Detection**: Bot messages and self-messages are automatically ignored
- **Case Insensitive**: Command names and aliases are matched case-insensitively
- **Permissions**: Permission checking is automatic when `RequiredUserPerms` or `RequiredBotPerms` are set

## Dependencies
- `github.com/bwmarrin/discordgo` - Discord API client
- `github.com/RhykerWells/durationutil` - Duration parsing utilities
- `github.com/sirupsen/logrus` - Logging

## License
See LICENSE file for details.
