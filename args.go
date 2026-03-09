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
	ValidateArg(arg *ParsedArg, data *Data) (any, bool)
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
	Int64	 = &Int64Arg{}
	User     = &UserArg{}
	Member   = &MemberArg{}
	Duration = &DurationArg{}
)

var (
	_ ArgumentType = (*StringArg)(nil)
	_ ArgumentType = (*IntArg)(nil)
	_ ArgumentType = (*Int64Arg)(nil)
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

func (s *StringArg) ValidateArg(arg *ParsedArg, data *Data) (any, bool) {
	v := arg.Raw

	if v == "" {
		return nil, false
	}

	if len(s.Options) > 0 {
		for _, option := range s.Options {
			if strings.EqualFold(v, option) {
				return v, true
			}
		}
		return nil, false
	}

	return v, true
}

type IntArg struct {
	Min int
	Max int
}

func (i *IntArg) Help() string {
	var maxStr string
	if i.Max != 0 {
		maxStr = fmt.Sprintf(" and below %d", i.Max)
	}
	return fmt.Sprintf("Whole number above %d%s", i.Min, maxStr)
}

func (i *IntArg) ValidateArg(arg *ParsedArg, data *Data) (any, bool) {
    v, err := strconv.Atoi(arg.Raw)
    if err != nil {
        return nil, false
    }

    if v < i.Min {
        return nil, false
    }

    if i.Max != 0 && v > i.Max {
        return nil, false
    }

    return v, true
}

type Int64Arg struct {
	Min int64
	Max int64
}

func (i *Int64Arg) Help() string {
	var maxStr string
	if i.Max != 0 {
		maxStr = fmt.Sprintf(" and below %d", i.Max)
	}
	return fmt.Sprintf("Whole number above %d%s", i.Min, maxStr)
}

func (i *Int64Arg) ValidateArg(arg *ParsedArg, data *Data) (any, bool) {
    v, err := strconv.ParseInt(arg.Raw, 10, 64)
    if err != nil {
        return nil, false
    }

    if v < i.Min {
        return nil, false
    }

    if i.Max != 0 && v > i.Max {
        return nil, false
    }

    return v, true
}

type UserArg struct{}

func (u *UserArg) Help() string {
	return "Mention/ID"
}

func (u *UserArg) ValidateArg(arg *ParsedArg, data *Data) (any, bool) {
	id := arg.Raw

	user, err := data.Session.User(id)
	if err != nil {
		return nil, false
	}

	return user, true
}

type MemberArg struct{}

func (m *MemberArg) Help() string {
	return "Mention/ID"
}

func (m *MemberArg) ValidateArg(arg *ParsedArg, data *Data) (any, bool) {
	id := arg.Raw

	member, err := data.Session.State.Member(data.Guild.ID, id)
	if err != nil {
		return nil, false
	}

	return member, true
}

type DurationArg struct{}

func (d *DurationArg) Help() string {
	return "Duration"
}

func (d *DurationArg) ValidateArg(arg *ParsedArg, data *Data) (any, bool) {
	v := arg.Raw

	duration, err := durationutil.ToDuration(v)
	if err != nil {
		return nil, false
	}

	return duration, true
}
