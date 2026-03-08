package dispatch

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/RhykerWells/durationutil"
)

// Arg defines the structure to pass argument data with
type Arg struct {
	Name string
	Type ArgumentType
}

type ArgumentType interface {
	ValidateArg(arg *ParsedArg, data *Data) bool
	Help() string
}

// ParsedArg represents a single argument parsed from a command invocation.
// Each argument contains the initial argument definition, the raw value passed to it, and the resolved value
type ParsedArg struct {
	Def   *Arg
	Raw   string
	Value any
}

var (
	String   = &StringArg{}
	Int      = &IntArg{}
	User     = &UserArg{}
	Member   = &MemberArg{}
	Duration = &DurationArg{}
)

var (
	_ ArgumentType = (*StringArg)(nil)
	_ ArgumentType = (*IntArg)(nil)
	_ ArgumentType = (*UserArg)(nil)
	_ ArgumentType = (*MemberArg)(nil)
	_ ArgumentType = (*DurationArg)(nil)
)

type StringArg struct {
	Options []string
}

func (s *StringArg) Help() string {
	if len(s.Options) > 0 {
		return fmt.Sprintf("%s", strings.Join(s.Options, "/"))
	}

	return "Text"
}

func (s *StringArg) ValidateArg(arg *ParsedArg, data *Data) bool {
	v := arg.Raw

	if v == "" {
		return false
	}

	if len(s.Options) > 0 {
		for _, option := range s.Options {
			if strings.EqualFold(v, option) {
				arg.Value = v
				return true
			}
		}
		return false
	}

	arg.Value = v
	return true
}

type IntArg struct {
	Min int64
	Max int64
}

func (i *IntArg) Help() string {
	var maxStr string
	if i.Max != 0 {
		maxStr = fmt.Sprintf(" and below %d", i.Max)
	}
	return fmt.Sprintf("Whole number above %d%s", i.Min, maxStr)
}

func (i *IntArg) ValidateArg(arg *ParsedArg, data *Data) bool {
    v, err := strconv.ParseInt(arg.Raw, 10, 64)
    if err != nil {
        return false
    }

    if v < i.Min {
        return false
    }

    if i.Max != 0 && v > i.Max {
        return false
    }

    arg.Value = v
    return true
}

type UserArg struct{}

func (u *UserArg) Help() string {
	return "Mention/ID"
}

func (u *UserArg) ValidateArg(arg *ParsedArg, data *Data) bool {
	id := arg.Raw

	user, err := data.Session.User(id)
	if err != nil {
		return false
	}

	arg.Value = user
	return true
}

type MemberArg struct{}

func (m *MemberArg) Help() string {
	return "Mention/ID"
}

func (m *MemberArg) ValidateArg(arg *ParsedArg, data *Data) bool {
	id := arg.Raw

	member, err := data.Session.State.Member(data.Guild.ID, id)
	if err != nil {
		return false
	}

	arg.Value = member
	return true
}

type DurationArg struct{}

func (d *DurationArg) Help() string {
	return "Duration"
}

func (d *DurationArg) ValidateArg(arg *ParsedArg, data *Data) bool {
	v := arg.Raw

	duration, err := durationutil.ToDuration(v)
	if err != nil {
		return false
	}

	arg.Value = duration
	return true
}
