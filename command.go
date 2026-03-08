package dispatch

// Command defines the general data that must be set during the addition of a new command
type Command struct {
	Command     string
	Category    CommandCategory
	Aliases     []string
	Description string

	Args           []*Arg
	ArgsRequired   int // Ignored if using combos
	ArgumentCombos [][]int

	RequiredUserPerms []int64
	RequiredBotPerms  []int64

	Run Run
}

// CommandCategory defines the available category types for commands
type CommandCategory struct {
	Name        string
	Description string
}

type Run func(data *Data) error

// RegisteredCommand defines the context required to access data surrounding a command
type RegisteredCommand struct {
	Command *Command
}
